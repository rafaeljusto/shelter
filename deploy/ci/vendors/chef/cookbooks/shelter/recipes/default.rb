#
# Cookbook Name:: shelter
# Recipe:: default
#
# Copyright 2014, Rafael Dantas Justo

include_recipe "mongodb::default"

shelter_latest = Chef::Config[:file_cache_path] + "/shelter_latest_amd64.deb"

remote_file shelter_latest do
  source "http://dl.bintray.com/rafaeljusto/deb/shelter_latest_amd64.deb"
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
    :base_path          => node["shelter"]["base_path"],
    :log_filename       => node["shelter"]["log_filename"],
    :db_name            => node["shelter"]["db_name"],
    :db_uri             => node["shelter"]["db_uri"],
    :scan_resolver      => node["shelter"]["scan_resolver"],
    :scan_resolver_port => node["shelter"]["scan_resolver_port"],
    :smtp_auth_user     => node["shelter"]["smtp_auth_user"],
    :smtp_auth_pwd      => node["shelter"]["smtp_auth_pwd"]
  )
end

execute "run-shelter" do
  command "service shelter start"
end