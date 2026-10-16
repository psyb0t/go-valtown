# Rewrites Val Town's OpenAPI document into a form oapi-codegen generates
# typed code for.
#
# 1. An `enum` of strings with no `type` gets `type: string`. Without a type,
#    oapi-codegen declares the field or parameter as interface{} and emits no
#    enum constants. Enums that mix value types are left alone.

def untyped_string_enum:
  type == "object"
  and has("enum")
  and (has("type") | not)
  and (.enum | length > 0)
  and all(.enum[]; type == "string");

walk(if untyped_string_enum then . + {type: "string"} else . end)
