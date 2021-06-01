// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ProjectArgs struct {
	ProjectID        string    `pulumi:"projectID"`
	ProjectName      *string   `pulumi:"projectName"`
	OrgID            *string   `pulumi:"orgID"`
	FolderID         *string   `pulumi:"folderID"`
	BillingAccountID *string   `pulumi:"billingAccountID"`
	ActivatedAPIs    *[]string `pulumi:"activatedApis"`
}

type Project struct {
	pulumi.ResourceState

	ProjectID pulumi.StringOutput `pulumi:"projectID"`
}

func NewProject(ctx *pulumi.Context,
	name string, args *ProjectArgs, opts ...pulumi.ResourceOption) (*Project, error) {
	if args == nil {
		args = &ProjectArgs{}
	}

	// Validation logic
	if args.OrgID != nil && args.FolderID != nil {
		return nil, fmt.Errorf("only one of `folderID` or `orgID` can be specified")
	}

	if args.OrgID == nil && args.FolderID == nil {
		return nil, fmt.Errorf("one of `folderID` or `orgID` must be specified")
	}

	// Construct Project Name
	var projectName string
	if args.ProjectName == nil {
		projectName = args.ProjectID
	} else {
		projectName = *args.ProjectName
	}

	component := &Project{}
	err := ctx.RegisterComponentResource(GCPProjectScaffold, name, component, opts...)
	if err != nil {
		return nil, err
	}

	projectCreationArgs := &organizations.ProjectArgs{
		ProjectId: pulumi.String(args.ProjectID),
		Name:      pulumi.String(projectName),
	}

	if args.BillingAccountID != nil && *args.BillingAccountID != "" {
		projectCreationArgs.BillingAccount = pulumi.String(*args.BillingAccountID)
	}

	if args.OrgID != nil {
		projectCreationArgs.OrgId = pulumi.String(*args.OrgID)
	}

	if args.FolderID != nil {
		projectCreationArgs.FolderId = pulumi.String(*args.FolderID)
	}

	project, err := organizations.NewProject(ctx, fmt.Sprintf("%s-%s", name, args.ProjectID),
		projectCreationArgs, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("error creating project: %v", err)
	}

	if args.ActivatedAPIs != nil {
		for _, api := range *args.ActivatedAPIs {
			serviceComponents := strings.Split(api, ".")
			_, err := projects.NewService(ctx, fmt.Sprintf("%s-%s-api", name, serviceComponents[0]),
				&projects.ServiceArgs{
					DisableDependentServices: pulumi.Bool(true),
					DisableOnDestroy:         pulumi.Bool(true),
					Project:                  project.ProjectId,
					Service:                  pulumi.String(api),
				}, pulumi.Parent(project))
			if err != nil {
				return nil, fmt.Errorf("unable to active service %q for project %q: %v", api,
					project.ProjectId, err)
			}
		}
	}

	component.ProjectID = project.ProjectId

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"projectID":   project.ProjectId,
		"projectName": pulumi.String(projectName),
	}); err != nil {
		return nil, err
	}

	return component, nil
}
