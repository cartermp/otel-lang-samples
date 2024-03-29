import {
    Context,
    context,
    propagation,
    trace,
    Tracer,
  } from '@opentelemetry/api';
  import express, { Express, Request, Response } from 'express';
  
  const app: Express = express();
  const hostname = '0.0.0.0';
  const port = 3000;
  
  // express supports async handlers but the @types definition is wrong: https://github.com/DefinitelyTyped/DefinitelyTyped/issues/50871
  // eslint-disable-next-line @typescript-eslint/no-misused-promises
  app.get('/', async (_req: Request, res: Response, next: any) => {
    try {
      res.statusCode = 200;
      res.setHeader('Content-Type', 'text/plain');
      const sayHello = () => 'Hello world!';
      const tracer: Tracer = trace.getTracer('hello-world-tracer');
      // new context based on current, with key/values added to baggage
      const ctx: Context = propagation.setBaggage(
        context.active(),
        propagation.createBaggage({
          for_the_children: { value: 'another important value' },
        }),
      );
      // within the new context, do some "work"
      await context.with(ctx, async () => {
        await tracer.startActiveSpan('sleep', async (span) => {
          console.log('saying hello to the world');
          span.setAttribute('message', 'hello-world');
          span.setAttribute('delay_ms', 100);
          await sleepy();
          console.log('sleeping a bit!');
          span.end();
        });
      });
      sayHello();
      res.end('Hello, World!\n');
    } catch (err) {
      next(err);
    }
  });
  
  function sleepy(): Promise<void> {
    return new Promise((resolve) => {
      setTimeout(() => {
        console.log('awake now!');
      }, 100);
      resolve();
    });
  }
  
  app.listen(port, hostname, () => {
    console.log(`Now listening on: http://${hostname}:${port}/`);
  });