'use strict';

import { HoneycombSDK } from '@honeycombio/opentelemetry-node';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';

// uses HONEYCOMB_API_KEY and OTEL_SERVICE_NAME environment variables
const sdk = new HoneycombSDK({
  instrumentations: [getNodeAutoInstrumentations()]
});

sdk.start()