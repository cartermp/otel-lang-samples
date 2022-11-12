require 'opentelemetry/sdk'
require 'opentelemetry/instrumentation/all'

begin
  OpenTelemetry::SDK.configure do |c|
    c.service_name = "ruby-sample"
    c.use_all()
  end
rescue OpenTelemetry::SDK::ConfigurationError => e
  puts "Error configuring Ruby OpenTelemetry SDK"
  puts e.inspect
end

Tracer = OpenTelemetry.tracer_provider.tracer("ruby-sample")