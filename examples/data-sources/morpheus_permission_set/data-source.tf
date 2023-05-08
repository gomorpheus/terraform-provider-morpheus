data "morpheus_permission_set" "source_one" {
  default_group_permission             = "full"
  default_cloud_permission             = "full"
  default_instance_type_permission     = "full"
  default_blueprint_permission         = "full"
  default_report_type_permission       = "full"
  default_persona_permission           = "full"
  default_catalog_item_type_permission = "full"
  default_vdi_pool_permission          = "full"
  default_workflow_permission          = "full"
  default_task_permission              = "full"

  feature_permission {
    name   = "ansible"
    access = "full"
  }

  group_permission {
    name   = "ansible"
    access = "full"
  }

  instance_type_permission {
    name   = "ansible"
    access = "full"
  }

  blueprint_permission {
    name   = "ansible"
    access = "full"
  }

  report_type_permission {
    name   = "ansible"
    access = "full"
  }

  persona_permission {
    name   = "ansible"
    access = "full"
  }

  catalog_item_type_permission {
    name   = "ansible"
    access = "full"
  }

  vdi_pool_permission {
    name   = "ansible"
    access = "full"
  }

  workflow_permission {
    name   = "ansible"
    access = "full"
  }

  task_permission {
    name   = "ansible"
    access = "full"
  }
}