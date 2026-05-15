# snow_table — reference module

A convenience wrapper around the three ServiceNow resources that together
define a custom table with custom columns: `servicenow_db_table`,
`servicenow_dictionary`, and `servicenow_choice`. Lets you declare a table
with all its columns (and choice values) in one block, while preserving
fine-grained Terraform diff handling at the column and choice level.

This is a **reference module**, not part of the provider binary. Copy it
into your own repository, or source it directly from this repository:

```hcl
module "demo" {
  source = "github.com/tylerhatton/terraform-provider-servicenow//examples/modules/snow_table?ref=v0.11.0"

  label = "TF Demo Table"

  columns = {
    u_environment = {
      column_label  = "Environment"
      internal_type = "choice"
      max_length    = 40
      choice        = 3 # dropdown
      mandatory     = true
      choices = {
        prod    = "Production"
        staging = "Staging"
        dev     = "Development"
      }
    }
    u_owner = {
      column_label  = "Owner"
      internal_type = "string"
      max_length    = 255
    }
    u_active = {
      column_label  = "Active"
      internal_type = "boolean"
      default_value = "true"
    }
  }
}

# The module's `name` output gives you the auto-generated table name
# (e.g. `u_tf_demo_table`). Use it on downstream resources:
output "demo_table_name" {
  value = module.demo.name
}
```

### Inserting rows into the new table

`servicenow_record` works against any table. To insert into the table
this module just created:

```hcl
resource "servicenow_record" "row_one" {
  table = module.demo.name
  fields = {
    u_environment = "prod"
    u_owner       = "platform-team"
    u_active      = "true"
  }
}
```

**Heads-up on ACLs.** ServiceNow auto-generates row-level ACLs for newly
created custom tables, and those ACLs often require a role (e.g. the
application's user role) that the API caller doesn't hold by default —
even `admin`. If `terraform apply` fails with `HTTP 403 Forbidden / User
Not Authorized` on the first row insert, either:

- Set `user_role = "<role-sys-id>"` on the module so a known role gates
  the table, then ensure your provider user holds that role, **or**
- Author the table's ACLs explicitly via `servicenow_acl` (note: ACL
  creation requires `security_admin` elevation), **or**
- Adjust the auto-generated ACLs in the ServiceNow UI after the table
  is created.

This is a ServiceNow configuration concern, not a provider/module bug —
the same 403 happens if you `curl` an insert against a freshly-created
custom table.

## Inputs

| Name | Type | Default | Description |
|---|---|---|---|
| `label` | string | — | Display label. ServiceNow generates the internal `name` from it. |
| `user_role` | string | `""` | Role sys_id required for end-user access. Empty = anyone. |
| `extendable` | bool | `false` | Allow other tables to extend this one. |
| `super_class` | string | `""` | sys_id of a parent table this one extends. |
| `columns` | map(object) | `{}` | Element name → column definition (see below). |

### Column object schema

Every column requires `column_label`. Every other attribute has a sensible default:

| Attribute | Type | Default | Notes |
|---|---|---|---|
| `column_label` | string | required | Human-readable label. |
| `internal_type` | string | `"string"` | `string`, `integer`, `boolean`, `choice`, `reference`, `glide_date`, etc. |
| `max_length` | number | `0` | Required for `string` columns. |
| `mandatory` | bool | `false` | |
| `read_only` | bool | `false` | |
| `active` | bool | `true` | |
| `display` | bool | `false` | If true, this is the table's display column. |
| `unique` | bool | `false` | |
| `default_value` | string | `""` | |
| `comments` | string | `""` | |
| `reference` | string | `""` | For `internal_type = "reference"`, the target table name. |
| `choice` | number | `0` | `0` none, `1` suggestion, `3` dropdown. |
| `choices` | map(string) | `{}` | For choice-typed columns: `{ value = label }`. Each entry becomes one `servicenow_choice` row. |

## Outputs

| Name | Type | Description |
|---|---|---|
| `id` | string | sys_id of the table. |
| `name` | string | Auto-generated table name (use as `table = module.x.name`). |
| `columns` | map(string) | Element name → dictionary sys_id. |
| `choices` | map(string) | `"<column>__<value>"` → choice sys_id. |

## Trade-offs vs raw resources

Use the module when you're defining a brand-new custom table and want a
single block per table. Drop down to the raw `servicenow_db_table` +
`servicenow_dictionary` + `servicenow_choice` resources when:

- You're attaching dictionary entries to a **built-in** table (e.g.
  adding `u_*` columns to `incident`) and the table itself isn't
  Terraform-managed.
- You need column-level metadata the module doesn't expose (e.g.
  `dependent`, `dependent_on_field`).
- You need to delete or rename columns out of order from how they were
  declared in the `columns` map.

Either approach can be mixed in the same configuration — the module's
outputs (`name`, `columns`, `choices`) are usable from outside the module.
