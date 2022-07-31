resource "morpheus_arm_spec_template" "tfexample_arm_spec_template_local" {
  name         = "tf-arm-spec-example-local"
  source_type  = "local"
  spec_content = <<TFEOF
{
  "type": "Microsoft.Storage/storageAccounts",
  "apiVersion": "2019-04-01",
  "name": "string",
  "location": "string",
  "tags": {
    "tagName1": "tagValue1",
    "tagName2": "tagValue2"
  },
  "sku": {
    "name": "string",
    "restrictions": [
      {
        "reasonCode": "string"
      }
    ]
  },
  "kind": "string",
  "identity": {
    "type": "SystemAssigned"
  },
  "properties": {
    "accessTier": "string",
    "allowBlobPublicAccess": "bool",
    "allowSharedKeyAccess": "bool",
    "azureFilesIdentityBasedAuthentication": {
      "activeDirectoryProperties": {
        "azureStorageSid": "string",
        "domainGuid": "string",
        "domainName": "string",
        "domainSid": "string",
        "forestName": "string",
        "netBiosDomainName": "string"
      },
      "directoryServiceOptions": "string"
    },
    "customDomain": {
      "name": "string",
      "useSubDomainName": "bool"
    },
    "encryption": {
      "keySource": "string",
      "keyvaultproperties": {
        "keyname": "string",
        "keyvaulturi": "string",
        "keyversion": "string"
      },
      "services": {
        "blob": {
          "enabled": "bool"
        },
        "file": {
          "enabled": "bool"
        }
      }
    },
    "isHnsEnabled": "bool",
    "largeFileSharesState": "string",
    "minimumTlsVersion": "string",
    "networkAcls": {
      "bypass": "string",
      "defaultAction": "string",
      "ipRules": [
        {
          "action": "Allow",
          "value": "string"
        }
      ],
      "virtualNetworkRules": [
        {
          "action": "Allow",
          "id": "string",
          "state": "string"
        }
      ]
    },
    "supportsHttpsTrafficOnly": "bool"
  }
}
TFEOF
}