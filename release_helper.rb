# frozen_string_literal: true

require 'yaml'

CHART_PATH = 'charts/uffizzi-controller/Chart.yaml'
CHART_VALUES_PATH = 'charts/uffizzi-controller/values.yaml'
CLUSTER_OPERATOR_DEPENDENCY_NAME = 'uffizzi-cluster-operator'
OPTION_REGEX = /^--/.freeze

def update_cluster_operator_version_command(options)
  chart = load_chart
  new_dependencies = chart['dependencies'].map do |d|
    d['version'] = options['version'] if d['name'] == CLUSTER_OPERATOR_DEPENDENCY_NAME
    d
  end

  chart['dependencies'] = new_dependencies
  save_chart(chart)
end

def update_controller_version_command(options)
  new_version = options['version']
  chart = load_chart
  chart['version'] = new_version
  chart['appVersion'] = new_version
  save_chart(chart)

  update_chart_values_image_tag_command(options)
end

def update_chart_values_image_tag_command(options)
  new_version = options['version']

  chart_values = load_chart_values
  repo = chart_values['image'].split(':').first
  new_image = [repo, new_version].join(':')
  chart_values['image'] = new_image
  save_chart_values(chart_values)
end

def load_chart
  YAML.safe_load(File.read(CHART_PATH))
end

def save_chart(chart)
  File.write(CHART_PATH, chart.to_yaml)
end

def load_chart_values
  YAML.safe_load(File.read(CHART_VALUES_PATH))
end

def save_chart_values(values)
  File.write(CHART_VALUES_PATH, values.to_yaml)
end

def parse_args
  options = {}
  commands = []
  prev_arg_option_name = nil

  ARGV.each do |arg|
    next commands << arg.strip if options.empty? && !OPTION_REGEX.match?(arg)

    if OPTION_REGEX.match?(arg)
      prev_arg_option_name = arg.gsub(OPTION_REGEX, '').strip
      options[prev_arg_option_name] = nil
    elsif !prev_arg_option_name.nil?
      options[prev_arg_option_name] = arg.strip
      prev_arg_option_name = nil
    end
  end

  { commands: commands, options: options }
end

def run
  args = parse_args
  command = args[:commands].first
  return if command.nil?

  options = args[:options] || {}
  send("#{command}_command", options)
end

run
