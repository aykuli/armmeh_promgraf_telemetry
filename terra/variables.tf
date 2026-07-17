# --- PROVIDERS ---
variable "cloud_id" {
  type = string
}
variable "folder_id" {
  type = string
}
variable "ssh_public_key" {
  type        = string
  description = "Ваш публичный SSH ключ"
}

variable "s3bucket" {
  type = string
}
# Zone	Subnet Name	CIDR Block (Example)	Purpose
# ru-central1-a	subnet-a	10.0.1.0/24	Resources in Zone A
# ru-central1-b	subnet-b	10.0.2.0/24	Resources in Zone B
# ru-central1-c	subnet-c	10.0.3.0/24	Resources in Zone C
variable "sg" {
  type = object({
    web = object({
      name        = string
      description = string
      labels      = map(string)
      ingress_rules = list(object({
        protocol    = string
        v4_cidr     = optional(list(string))
        port        = optional(number)
        from_port   = optional(number)
        to_port     = optional(number)
        description = optional(string)
      }))
      egress_rules = optional(list(object({
        protocol    = string
        v4_cidr     = optional(list(string))
        port        = optional(number)
        from_port   = optional(number)
        to_port     = optional(number)
        description = optional(string)
      })))
    })
  })
}


