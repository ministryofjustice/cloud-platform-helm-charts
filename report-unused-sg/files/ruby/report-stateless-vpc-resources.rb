#!/usr/bin/env ruby

require "pry-byebug"
require "json"
require "bundler/setup"
require "aws-sdk-ec2"

def nat_gateway_ids_for_vpc(client, vpc_id)
  filter = [ { name: "vpc-id", values: [vpc_id] } ]
  data = client.describe_nat_gateways(filter: filter)
  data.nat_gateways.map { |ng| ng.nat_gateway_id }.sort
end

def subnets_ids_for_vpc(client, vpc_id)
  filter = [ { name: "vpc-id", values: [vpc_id] } ]
  data = client.describe_subnets(filter: filter)
  data.subnets.map { |sn| sn.subnet_id }.sort
end

def route_tables_for_vpc(client, subnet_id)
  filter = [ { name: "subnet_id", values: [subnet_id] } ]
  data = client.describe_route_tables(filter: filter)
  data.associations.map { |rt| rt.route_table_id }.sort
end

def internet_gateway_ids_for_vpc(client)
    data = client.describe_internet_gateways()
    data.internet_gateways.map { |ng| ng.internet_gateway_id }.sort
end

def vpc_ids(client)
    data = client.describe_vpcs()
    data.vpcs.map { |ng| ng.vpc_id }.sort
end


def nat_gateway_ids_from_terraform_state(statefile)
  str = File.read("terraform.tfstate")
  data = JSON.parse(str)
  list = data["resources"]
  nat_gateway = list.filter { |m| m["name"] == "private_nat_gateway" }.first
  nat_gateway["instances"].map { |ng| ng.dig("attributes", "nat_gateway_id") }.sort
end


def internet_gateway_ids_from_terraform_state(statefile)
    str = File.read("terraform.tfstate")
    data = JSON.parse(str)
    list = data["resources"]
    nat_gateway = list.filter { |m| m["name"] == "public_internet_gateway" }.first
    nat_gateway["instances"].map { |ng| ng.dig("attributes", "gateway_id") }.sort
end

#binding.pry

ec2 = Aws::EC2::Client.new(region:'eu-west-2', profile: ENV["AWS_PROFILE"])

#******** Get the natgateway ids **********************

#pp nat_gateway_ids_for_vpc(ec2, "vpc-0c16457fd570a1f0b")
#pp nat_gateway_ids_from_terraform_state("terraform.tfstate")

#********Get the internet gateway*************************************************

#pp internet_gateway_ids_for_vpc(ec2)
#pp internet_gateway_ids_from_terraform_state("terraform.tfstate")

#********Get the vpc ids***********************
#pp vpc_ids(ec2)

#*********Get the subnets*******************

pp subnets_ids_for_vpc(ec2, "vpc-0922049e634e2dc7b")

#pp route_tables_for_vpc(ec2, "subnet-05457126ff452defb")



