require "opentelemetry/sdk"

class AppController < ApplicationController
  def index
    current_span = OpenTelemetry::Trace.current_span
    current_span.set_attribute("banana", "butt")

    Tracer.in_span("app", attributes: { "hello" => "world", "some.number" => 1024 }) do |span|
      span.add_attributes(
        {
          OpenTelemetry::SemanticConventions::Trace::HTTP_METHOD => "GET",
          OpenTelemetry::SemanticConventions::Trace::HTTP_URL => "https://opentelemetry.io/",
        })
      child_work(current_span.context)
    end

    current_span.add_event("poopy-butt", attributes: { "poot" => "toot" })
  end

  def child_work(span_context)
    link = OpenTelemetry::Trace::Link.new(span_context, attributes: { "linky" => "loo" })
    Tracer.in_span("child", links: [link]) do |span|
      span.add_attributes(
        {
          "my.cool.attribute" => "a value",
          "my.first.name" => "Oscar"
        })

      begin
        1/0
      rescue Exception => ex
        span.status = OpenTelemetry::Trace::Status.error(ex.message)
        span.record_exception(ex, attributes: { "some.attribute" => 12 })
      end
    end
  end
end
