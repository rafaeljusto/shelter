#
# Cookbook Name:: shelter
# Recipe:: default
#
# Copyright 2014, YOUR_COMPANY_NAME
#
# All rights reserved - Do Not Redistribute
#

include_recipe "mongodb::default"

shelter_latest = Chef::Config[:file_cache_path] + "/rafaeljusto/deb/shelter/_latestVersio"

remote_file shelter_latest do
  source "https://bintray.com/rafaeljusto/deb/shelter/_latestVersion"
  mode "0644"
end

execute "install-shelter" do
  command "dpkg --no-triggers -i " + shelter_latest
  creates "/usr/shelter/etc/shelter.conf"
end

template "/usr/shelter/etc/shelter.conf" do
  source "shelter.conf.chef.sample"
  mode 0755
  owner "root"
  group "root"
  variables(
    
  )
end
