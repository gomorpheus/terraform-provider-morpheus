## Unreleased

NOTES:

* The `morpheus-sdk` dependcy has been upgraded to version 0.2.10.
* Add label support for additional Morpheus resources.
* Update inputs to support additional configuration parameters (i.e. - editable, verify pattern, etc).

FEATURES:

* **New Data Source:** `morpheus_vro_workflow`
* **New Resource:** `morpheus_active_directory_identity_source`
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
