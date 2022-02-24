## Unreleased


## 0.4.0 (February 24, 2022)

NOTES:

* The `morpheus_tenant` resource has been updated to fix an invalid api call that prevented the creation of tenants using the provider.

* The name and data type of the `base_role` attribute for the `morpheus_tenant` resource has been changed. The new name is `base_role_id` and the data type is an integer instead of a string.

* Source header support for REST API option lists has been added to the `morpheus_rest_option_list` resource.

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
