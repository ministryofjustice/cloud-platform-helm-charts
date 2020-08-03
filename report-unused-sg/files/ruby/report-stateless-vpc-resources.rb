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
  filters = [ { name: "vpc-id", values: [vpc_id] } ]
  data = client.describe_subnets(filters: filters)
  data.subnets.map { |sn| sn.subnet_id }.sort
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

def route_tables_for_subnet(client, subnet_id)
  filters = [ { name: "association.subnet-id", values: [subnet_id] } ]
  data = client.describe_route_tables(filters: filters)
  data.route_tables.map { |rt| rt.route_table_id }.sort
end

def route_table_associations_for_subnet(client, subnet_id)
  filters = [ { name: "association.subnet-id", values: [subnet_id] } ]
  data = client.describe_route_tables(filters: filters)
  puts data
end


#def route_table_associations_for_subnet(client, route_table_id)
#  filters = [ { name: "association.route-id", values: [route_table_id] } ]
 # data = client.describe_route_tables({route_table_ids: [route_table_id, ],})
 # data = client.describe_route_tables(filters: filters)
 # puts data
 # puts data[:route_table_id]
#end


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

#pp subnets_ids_for_vpc(ec2, "vpc-0c16457fd570a1f0b")


#****Get the route table ids for each subnet*******

#pp route_tables_for_subnet(ec2, "subnet-0386053fb92e69994")
#pp route_tables_for_subnet(ec2, "subnet-05fe20e709ed69201")


pp route_table_associations_for_subnet(ec2, "subnet-05fe20e709ed69201")