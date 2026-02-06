import { z } from "zod";

import { createTRPCRouter, publicProcedure } from "~/server/api/trpc";
import { gateway } from "~/server/api/gateway";
import { CardDataSchema } from "../../../../../api/contracts";

/**
 * Card router — proxies requests to the Go API Gateway.
 *
 * All card operations go through the gateway which routes to the
 * appropriate microservice (card-generator, image-renderer, importer).
 */
export const cardRouter = createTRPCRouter({
  /** Check API Gateway + service health */
  health: publicProcedure.query(async () => {
    return gateway.health();
  }),

  /** Generate a new card */
  generate: publicProcedure
    .input(CardDataSchema)
    .mutation(async ({ input }) => {
      return gateway.generateCard(input);
    }),

  /** Get a card by ID */
  getById: publicProcedure
    .input(z.object({ id: z.string().uuid() }))
    .query(async ({ input }) => {
      return gateway.getCard(input.id);
    }),

  /** Analyze a card for tags and synergies */
  analyze: publicProcedure
    .input(CardDataSchema)
    .mutation(async ({ input }) => {
      return gateway.analyzeCard(input);
    }),

  /** Get import job status */
  importStatus: publicProcedure
    .input(z.object({ jobId: z.string().uuid() }))
    .query(async ({ input }) => {
      return gateway.getImportStatus(input.jobId);
    }),

  /** Clear all cards (admin) */
  clearAll: publicProcedure.mutation(async () => {
    return gateway.clearCards();
  }),
});
