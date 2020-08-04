#!/usr/bin/env ruby

require "pry-byebug"
require "json"
require "bundler/setup"
require "aws-sdk-ec2"
require "aws-sdk-s3"

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
  data[:route_tables][0][:associations].each do |route|
  #pp route[:route_table_association_id]
    pp route[:route_table_id]
  end
end


def get_state_vpc_names()

  vpc_state_names = [] 
  
  # Create a file containing the workspace list (vpc names)
  puts `terraform init`
  File.open("vpc_state_names.txt","w") do |f|
    f.puts `terraform workspace list`
  end

  # Populate the vpc array with the vpc names from the file 
  File.open('vpc_state_names.txt').each { |vpc| vpc_state_names << vpc.delete(' ') }

  # Remove the first element as this is the default vpc
  vpc_state_names.shift(1)

  return vpc_state_names

end 


def get_vpc_ids_from_aws(client)

  vpc_actual_ids = []

  # Get the whole data for all vpcs 
  data = client.describe_vpcs()

  # Create a file containing the vpc ids only
  File.open("vpc_actual_ids.txt","w") do |f|
    f.puts data.vpcs.map { |vpc| vpc.vpc_id }.sort
  end

  File.open('vpc_actual_ids.txt').each { |vpc| vpc_actual_ids << vpc }

  #vpc_actual_ids.shift(1)

  return vpc_actual_ids

end 

def get_vpc_ids_from_state(statefile)
  str = File.read(statefile)
  data = JSON.parse(str)
  vpc_id = data['outputs']['vpc_id']['value']
  return vpc_id
end


def download_network_state(s3)

  # Download the network state file from s3 for a given vpc name 
  obj = s3.bucket('cloud-platform-terraform-state').object('cloud-platform-network/iawan-test/terraform.tfstate')
  obj.get(response_target: 'state.tfstate')
end

s3 = Aws::S3::Resource.new(region:'eu-west-1', profile: ENV["AWS_PROFILE"])

#binding.pry
ec2 = Aws::EC2::Client.new(region:'eu-west-2', profile: ENV["AWS_PROFILE"])


download_network_state(s3)

pp get_vpc_ids_from_state('state.tfstate')


#Get the VPC names in the state
#vpc_names_in_state = get_state_vpc_names()
#vpc_names_in_state.each do |vpc_name|
#  puts vpc_name
#end

#Get the actual VPC ids
vpc_ids_aws = get_vpc_ids_from_aws(ec2)
vpc_ids_aws.each do |vpc_id|
  puts vpc_id
end




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
#route_tables_for_subnet(ec2, "subnet-05fe20e709ed69201")


