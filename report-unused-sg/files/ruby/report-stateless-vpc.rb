#!/usr/bin/env ruby

require "pry-byebug"
require "json"
require "bundler/setup"
require "aws-sdk-ec2"
require "aws-sdk-s3"

def list_vpc_from_s3_state_folder(s3)
  vpc_names = []
  obj = s3.bucket("cloud-platform-terraform-state").objects(prefix:'cloud-platform-network', delimiter: '').collect(&:key) 
  obj.each do | vpc_name |
    vpc_names.push vpc_name
  end
  vpc_names.each { |item| puts item }
end

#binding.pry
s3 = Aws::S3::Resource.new(region:'eu-west-1', profile: ENV["AWS_PROFILE"])

pp list_vpc_from_s3_state_folder(s3)


