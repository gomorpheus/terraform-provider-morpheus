## UNRELEASED

FEATURES:

* **New Resource:** `morpheus_vsphere_cloud_datastore_configuration`

## 0.9.6 (November 14, 2023)

NOTES:

* Add initial support for properly reconciling missing resources that are managed with Terraform.
* Updated the `morpheus-go-sdk` library from v0.3.4 to v0.3.6
* Add support for labels on `morpheus_node_type`, `morpheus_instance_layout`, and `morpheus_instance_type` resources.
* Updated authentication documentation to properly document the use of the subtenant domain setting when using environment variables when authenticating.
* Updated the `morpheus_ansible_tower_integration` resource to fix an issue with the password continually indicating a change.
* Added support for utilizing credentials to the `morpheus_ansible_tower_integration` resource.

FEATURES:

* **New Resource:** `morpheus_credential`
* **New Resource:** `morpheus_cypher_secret`
* **New Resource:** `morpheus_cypher_tfvars`
* **New Resource:** `morpheus_library_script_task`
* **New Resource:** `morpheus_library_template_task`

## 0.9.5 (October 6, 2023)

NOTES:

* Updated the `terraform-plugin-docs` library from v0.14.1 to v0.16.0
* Updated the `terraform-plugin-sdk` library from v2.25.0 to v2.29.0
* Added a `version` parameter for the `instance_layout` data source to properly handle multiple instance layouts with the same name.
* Added the `tenant_subdomain` setting for the provider configuration to properly support authenticating with username and password to a subetnant. [#74](https://github.com/gomorpheus/terraform-provider-morpheus/issues/74)
* Fixed the issue that was causing the logo/image path to continually indicate that there is a change despite nothing being changed. [#146](https://github.com/gomorpheus/terraform-provider-morpheus/issues/146)

FEATURES:

* **New Data Source:** `morpheus_domain`
* **New Resource:** `morpheus_azure_cloud`

## 0.9.4 (August 14, 2023)

NOTES:

* Add support for managing a SAML identity source. [#140](https://github.com/gomorpheus/terraform-provider-morpheus/issues/140) 
* Add support for managing user roles. [#114](https://github.com/gomorpheus/terraform-provider-morpheus/issues/114) 

FEATURES:

* **New Data Source:** `morpheus_catalog_item_type`
* **New Data Source:** `morpheus_permission_set`
* **New Data Source:** `morpheus_security_package`
* **New Data Source:** `morpheus_vdi_pool`
* **New Resource:** `morpheus_ipv4_ip_pool`
* **New Resource:** `morpheus_saml_identity_source`
* **New Resource:** `morpheus_security_package`
* **New Resource:** `morpheus_servicenow_integration`
* **New Resource:** `morpheus_user_role`

## 0.9.3 (June 13, 2023)

NOTES:

* Updated the `aws_cloud` resource to add support for using host IAM credentials when authenticating to the cloud. [#103](https://github.com/gomorpheus/terraform-provider-morpheus/issues/103) 
* Updated the `api_option_list`, `manual_option_list`, and `rest_option_list` resources to better handle the difference in the payload returned from the API and the payload defined by Terraform. The payloads are now being compared after a trim operation has been performed on the payload passed by Terraform to address cases in which a HEREDOC is used that includes additional spacing for readability. [#128](https://github.com/gomorpheus/terraform-provider-morpheus/issues/128)
* Updated the `vsphere_cloud` resource to support importing existing VMware vSphere cloud integrations. [#129](https://github.com/gomorpheus/terraform-provider-morpheus/issues/129)
* Updated the logic for setting the state for the `provisioning_workflow` resource to properly account for the API returning the tasks in API versions prior to 5.5.x in an out of order sequence. This resulted in an inconsistent state and plans constantly indicating that there were changes to be made despite the real configuration not chaning. [#116](https://github.com/gomorpheus/terraform-provider-morpheus/issues/116)

FEATURES:

* **New Data Source:** `morpheus_git_integration`
* **New Data Source:** `morpheus_servicenow_workflow`
* **New Resource:** `morpheus_workflow_job`
* **New Resource:** `morpheus_provision_approval_policy`
* **New Resource:** `morpheus_delete_approval_policy`

## 0.9.2 (June 5, 2023)

NOTES:

* Update the `shell_script_task` resource to support local repository references and visibility attributes.
* Update the evaluation logic for sending the user group id data payload when creating an instance using the `vsphere_instance` resource. The logic previously caused an error when the user_group_id attribute was not set despite it being an optional attribute. [#121](https://github.com/gomorpheus/terraform-provider-morpheus/issues/121)
* Update the `vsphere_cloud` resource to add support for credentials referenced from the credential store. [#120](https://github.com/gomorpheus/terraform-provider-morpheus/issues/120)

## 0.9.1 (May 12, 2023)

NOTES:

* The `morpheus-sdk` dependcy has been upgraded to version 0.3.3.
* Updated the `morpheus_task_job` resource to properly read all object attributes. [#113](https://github.com/gomorpheus/terraform-provider-morpheus/issues/113)
* Updated the `morpheus_task_job` resource to add support for labels and support the new dynamic automation targeting feature in which instance or server labels can be used for the target selection.
* Add label support for automation task and workflow resources (i.e - provisioning workflow, ansible playbook task, python script task, etc).
* Add label support for option list resources.

FEATURES:

* **New Resource:** `morpheus_resource_pool_group`
* **New Resource:** `morpheus_license`
* **New Resource:** `morpheus_key_pair`

## 0.9.0 (April 26, 2023)

NOTES:

* The `morpheus-sdk` dependcy has been upgraded to version 0.3.2.
* The `terraform-plugin-docs` dependcy has been upgraded to version 0.14.1.
* The `terraform-plugin-sdk` dependcy has been upgraded to version 2.25.0.
* Add label support for additional Morpheus resources.
* Update inputs to support additional configuration parameters (i.e. - editable, verify pattern, etc).
* Updated provisioning workflow resource phase documentation and added validation support [#96](https://github.com/gomorpheus/terraform-provider-morpheus/issues/96)
* Fixed a bug with the provisioning workflow resource not properly reading tasks, which impacted updates and state import operations [#96](https://github.com/gomorpheus/terraform-provider-morpheus/issues/96)
* Updated the `morpheus_email_task` resource to support repository and url source types. [#97](https://github.com/gomorpheus/terraform-provider-morpheus/issues/97)
* Updated the `morpheus_powershell_task` resource documentation and fixed an issue with the execute_target attribute not properly being set on import. [#98](https://github.com/gomorpheus/terraform-provider-morpheus/issues/98)
* Updated the `morpheus_node_type` resource to remove the computed attribute for the `extra_options` attribute that was causing issues for non-vsphere resources following a resource import. [#100](https://github.com/gomorpheus/terraform-provider-morpheus/issues/100)
* Updated the `morpheus_instance_catalog_item` resource to properly set the state for config and visibility during resource import. [#102](https://github.com/gomorpheus/terraform-provider-morpheus/issues/102)

FEATURES:

* **New Data Source:** `morpheus_ansible_tower_inventory`
* **New Data Source:** `morpheus_ansible_tower_job_template`
* **New Data Source:** `morpheus_vro_workflow`
* **New Resource:** `morpheus_active_directory_identity_source`
* **New Resource:** `morpheus_ansible_tower_task`
* **New Resource:** `morpheus_guidance_setting`
* **New Resource:** `morpheus_monitoring_setting`
* **New Resource:** `morpheus_vro_integration`
* **New Resource:** `morpheus_vro_task`

## 0.8.0 (February 23, 2023)

NOTES:

* The `morpheus-sdk` dependcy has been upgraded to version 0.2.9.
* Fix the `morpheus_provisioning_workflow` resource to properly support the "all" platform setting options.
* Add label support for various Morpheus resources.
* Add custom option support for the `vsphere_instance` resource.
* Fix issue #58 - Incorrect default monitoring check for node types.
* Remove unnecessary request logging

FEATURES:

* **New Data Source:** `morpheus_network_group`
* **New Resource:** `morpheus_api_option_list`
* **New Resource:** `morpheus_app_blueprint_catalog_item`
* **New Resource:** `morpheus_aws_cloud`
* **New Resource:** `morpheus_radio_list_option_type`
* **New Resource:** `morpheus_textarea_option_type`

## 0.7.0 (December 16, 2022)

NOTES:

* The `morpheus-sdk` dependcy has been upgraded to version 0.2.5.
* Fix the `morpheus_ansible_integration` resource to properly set the default branch for the integration.

FEATURES:

* **New Data Source:** `morpheus_storage_bucket`
* **New Data Source:** `morpheus_user_group`
* **New Resource:** `morpheus_ansible_tower_integration`
* **New Resource:** `morpheus_backup_setting`
* **New Resource:** `morpheus_boot_script`
* **New Resource:** `morpheus_cypher_access_policy`
* **New Resource:** `morpheus_delayed_delete_policy`
* **New Resource:** `morpheus_instance_catalog_item`
* **New Resource:** `morpheus_motd_policy`
* **New Resource:** `morpheus_power_schedule_policy`
* **New Resource:** `morpheus_preseed_script`
* **New Resource:** `morpheus_puppet_integration`
* **New Resource:** `morpheus_tag_policy`
* **New Resource:** `morpheus_user_group_creation_policy`
* Add `repository_ids` attribute to the `git_integration` resource for using the repository ID for git based integration references such as shell script automation tasks.
* Add support for defining the assigned tenants for policy resources (i.e. - backup creation, budget, cypher access, etc.)

## 0.6.0 (September 12, 2022)

NOTES:

* The `morpheus-sdk` dependcy has been upgraded to version 0.1.8.
* Fix retry default settings for automation task resources.

FEATURES:

* **New Data Source:** `morpheus_blueprint`
* **New Data Source:** `morpheus_budget`
* **New Data Source:** `morpheus_cluster_type`
* **New Data Source:** `morpheus_credential`
* **New Data Source:** `morpheus_file_template`
* **New Data Source:** `morpheus_job`
* **New Data Source:** `morpheus_node_type`
* **New Data Source:** `morpheus_option_list`
* **New Data Source:** `morpheus_policy`
* **New Data Source:** `morpheus_power_schedule`
* **New Data Source:** `morpheus_provision_type`
* **New Data Source:** `morpheus_script_template`
* **New Data Source:** `morpheus_virtual_image`
* **New Resource:** `morpheus_ansible_integration`
* **New Resource:** `morpheus_cluster_layout`
* **New Resource:** `morpheus_docker_registry_integration`
* **New Resource:** `morpheus_file_template`
* **New Resource:** `morpheus_git_integration`
* **New Resource:** `morpheus_instance_layout`
* **New Resource:** `morpheus_instance_type`
* **New Resource:** `morpheus_node_type`
* **New Resource:** `morpheus_scale_threshold`
* **New Resource:** `morpheus_script_template`

## 0.5.0 (July 31, 2022)

NOTES:

* The required Golang version has been changed from 1.14 to 1.17 to support the recent versions of the Terraform plugin sdk.

* The `terraform-plugin-docs` dependency has been upgraded to 0.13.0.

* The `terraform-plugin-sdk` dependcy has been upgraded to version 2.18.0.

* The `morpheus_vsphere_cloud` resource has been updated to support properly managing the user credentials. This was enabled due to an API change to support proper credential handling via checksum comparisons between the Terraform value and the checksummed value returned by the API.

* Update existing data sources to support using the id of the Morpheus object in addition to the name of the object.

FEATURES:

* **New Data Source:** `morpheus_integration`
* **New Data Source:** `morpheus_price`
* **New Data Source:** `morpheus_price_set`
* **New Data Source:** `morpheus_tenant`
* **New Data Source:** `morpheus_spec_template`
* **New Resource:** `morpheus_arm_app_blueprint`
* **New Resource:** `morpheus_arm_spec_template`
* **New Resource:** `morpheus_backup_creation_policy`
* **New Resource:** `morpheus_budget_policy`
* **New Resource:** `morpheus_cloud_formation_app_blueprint`
* **New Resource:** `morpheus_cloud_formation_spec_template`
* **New Resource:** `morpheus_cluster_resource_name_policy`
* **New Resource:** `morpheus_groovy_script_task`
* **New Resource:** `morpheus_helm_app_blueprint`
* **New Resource:** `morpheus_helm_spec_template`
* **New Resource:** `morpheus_hostname_policy`
* **New Resource:** `morpheus_instance_name_policy`
* **New Resource:** `morpheus_javascript_task`
* **New Resource:** `morpheus_kubernetes_app_blueprint`
* **New Resource:** `morpheus_kubernetes_spec_template`
* **New Resource:** `morpheus_max_containers_policy`
* **New Resource:** `morpheus_max_memory_policy`
* **New Resource:** `morpheus_max_storage_policy`
* **New Resource:** `morpheus_network_quota_policy`
* **New Resource:** `morpheus_powershell_script_task`
* **New Resource:** `morpheus_price`
* **New Resource:** `morpheus_price_set`
* **New Resource:** `morpheus_restart_task`
* **New Resource:** `morpheus_router_quota_policy`
* **New Resource:** `morpheus_ruby_script_task`
* **New Resource:** `morpheus_service_plan`
* **New Resource:** `morpheus_shell_script_task`
* **New Resource:** `morpheus_terraform_app_blueprint`
* **New Resource:** `morpheus_user_creation_policy`
* **New Resource:** `morpheus_wiki_page`
* **New Resource:** `morpheus_workflow_catalog_item`
* **New Resource:** `morpheus_write_attributes_task`

## 0.4.0 (February 24, 2022)

NOTES:

* The `morpheus_tenant` resource has been updated to fix an invalid api call that prevented the creation of tenants using the provider.

* The name and data type of the `base_role` attribute for the `morpheus_tenant` resource has been changed. The new name is `base_role_id` and the data type is an integer instead of a string.

* Source header support for REST API option lists has been added to the `morpheus_rest_option_list` resource.

* Update the reference to the morpheus-go-sdk to use a tagged version to support the automated release process.

FEATURES:

* **New Data Source:** `morpheus_contact`
* **New Data Source:** `morpheus_execute_schedule`
* **New Data Source:** `morpheus_tenant_role`
* **New Data Source:** `morpheus_workflow`
* **New Resource:** `morpheus_contact`
* **New Resource:** `morpheus_execute_schedule`
* **New Resource:** `morpheus_max_cores_policy`
* **New Resource:** `morpheus_max_hosts_policy`
* **New Resource:** `morpheus_max_vms_policy`
* **New Resource:** `morpheus_task_job`
* **New Resource:** `morpheus_workflow_policy`

## 0.3.1 (September 23, 2021)

NOTES:

* Documentation updates to the README and installation instructions for Windows.

## 0.3.0 (August 18, 2021)

NOTES:

* Migration of the provider versioning to include the patch number in the versioning.

FEATURES:

* **New Resource:** `morpheus_terraform_spec_template`
* **New Resource:** `morpheus_python_script_task`
* **New Resource:** `morpheus_ansible_playbook_task`

## 0.2 (May 17, 2021)

FEATURES:

* **New Data Source:** `morpheus_cloud`
* **New Data Source:** `morpheus_environment`
* **New Data Source:** `morpheus_group`
* **New Data Source:** `morpheus_instance_layout`
* **New Data Source:** `morpheus_instance_type`
* **New Data Source:** `morpheus_network`
* **New Data Source:** `morpheus_option_type`
* **New Data Source:** `morpheus_plan`
* **New Data Source:** `morpheus_resource_pool`
* **New Data Source:** `morpheus_task`
* **New Resource:** `morpheus_checkbox_option_type`
* **New Resource:** `morpheus_hidden_option_type`
* **New Resource:** `morpheus_manual_option_list`
* **New Resource:** `morpheus_number_option_type`
* **New Resource:** `morpheus_operational_workflow`
* **New Resource:** `morpheus_password_option_type`
* **New Resource:** `morpheus_provisioning_workflow`
* **New Resource:** `morpheus_rest_option_list`
* **New Resource:** `morpheus_select_list_option_type`
* **New Resource:** `morpheus_tenant`
* **New Resource:** `morpheus_text_option_type`
* **New Resource:** `morpheus_typeahead_option_type`
* **New Resource:** `morpheus_vsphere_cloud`
* **New Resource:** `morpheus_vsphere_instance`

## 0.1 (November 27, 2019)

NOTES:

* This is a **BETA** version of the Morpheus Terraform Provider.

FEATURES:

* **New Resource:** `morpheus_cloud`
* **New Resource:** `morpheus_group`
* **New Resource:** `morpheus_instance`
* **New Resource:** `morpheus_network_domain`
