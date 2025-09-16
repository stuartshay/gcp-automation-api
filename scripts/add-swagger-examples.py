#!/usr/bin/env python3
"""
Script to add Swagger 2.0 request body examples to swagger.json
This adds Basic and Advanced examples by modifying parameter definitions
"""

import json
import sys


def add_swagger_examples(swagger_spec):
    """Add examples to Swagger 2.0 parameter definitions"""

    # Define examples for each model
    model_examples = {
        "BucketRequest": {
            "Basic": {
                "name": "my-simple-storage-bucket",
                "location": "us-central1",
                "storage_class": "STANDARD",
            },
            "Advanced": {
                "name": "my-enterprise-bucket-2024",
                "location": "us-central1",
                "storage_class": "STANDARD",
                "labels": {
                    "environment": "production",
                    "team": "platform",
                    "cost-center": "engineering",
                },
                "versioning": True,
                "kms_key_name": "projects/velvety-byway-327718/locations/us-central1/keyRings/bucket-encryption/cryptoKeys/bucket-key",
                "retention_policy": {
                    "retention_period_seconds": 7776000,
                    "is_locked": False,
                },
                "uniform_bucket_level_access": True,
                "public_access_prevention": "enforced",
            },
        },
        "ProjectRequest": {
            "Basic": {
                "project_id": "my-simple-project-2024",
                "display_name": "My Simple Project",
            },
            "Advanced": {
                "project_id": "enterprise-app-prod-2024",
                "display_name": "Enterprise Application - Production",
                "parent_id": "123456789012",
                "parent_type": "organization",
                "labels": {
                    "environment": "production",
                    "team": "backend",
                    "cost-center": "engineering",
                    "compliance": "sox",
                },
            },
        },
        "FolderRequest": {
            "Basic": {
                "display_name": "Development Environment",
                "parent_id": "123456789012",
                "parent_type": "organization",
            },
            "Advanced": {
                "display_name": "Production - North America Region",
                "parent_id": "987654321098",
                "parent_type": "folder",
            },
        },
    }

    # Add examples to definitions
    if "definitions" in swagger_spec:
        for model_name, examples in model_examples.items():
            full_model_name = f"models.{model_name}"
            if full_model_name in swagger_spec["definitions"]:
                definition = swagger_spec["definitions"][full_model_name]

                # Add x-examples extension for Swagger UI
                definition["x-examples"] = {
                    "Basic": {
                        "summary": "Basic Example",
                        "description": f"Simple {model_name.lower().replace('request', '')} with minimal required fields",
                        "value": examples["Basic"],
                    },
                    "Advanced": {
                        "summary": "Advanced Example",
                        "description": f"Enterprise {model_name.lower().replace('request', '')} with all available options",
                        "value": examples["Advanced"],
                    },
                }

                print(f"Added examples to {full_model_name}")

    return swagger_spec


def main():
    swagger_file = "docs/swagger.json"

    try:
        # Read swagger.json
        with open(swagger_file, "r") as f:
            swagger_spec = json.load(f)

        # Add examples
        updated_spec = add_swagger_examples(swagger_spec)

        # Write back to file
        with open(swagger_file, "w") as f:
            json.dump(updated_spec, f, indent=2)

        print(f"Successfully added Swagger 2.0 examples to {swagger_file}")

    except Exception as e:
        print(f"Error processing swagger file: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
