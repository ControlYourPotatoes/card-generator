import { z } from "zod";

// Base Card Data Schema (shared across endpoints)
export const CardDataSchema = z.object({
  name: z.string().min(1, "Name required"),
  cost: z.number().int().min(0).max(99),
  card_type: z.enum(["creature", "spell", "artifact", "incantation", "anthem"]),
  effect: z.string().min(1),
  keywords: z.array(z.string()).default([]),
  metadata: z.record(z.string()).optional(),
});

// API v1 Contracts with Zod validation

// POST /api/v1/cards/generate - Generate card from JSON data
export const GenerateCardRequestSchema = CardDataSchema;
export type GenerateCardRequest = z.infer<typeof GenerateCardRequestSchema>;

export const GenerateCardResponseSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  card_type: z.string(),
  tags: z.array(z.string()),
  metadata: z.record(z.string()),
  imageUrl: z.string().url().optional(),
  status: z.enum(["success", "processing", "failed"]),
});
export type GenerateCardResponse = z.infer<typeof GenerateCardResponseSchema>;

// GET /api/v1/cards/:id - Retrieve card metadata
export const GetCardRequestSchema = z.object({
  params: z.object({ id: z.string().uuid() }),
});
export type GetCardRequest = z.infer<typeof GetCardRequestSchema>;

export const CardResponseSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  cost: z.number(),
  card_type: z.string(),
  effect: z.string(),
  tags: z.array(z.string()),
  metadata: z.record(z.string()),
  imageUrl: z.string().url().optional(),
});
export type CardResponse = z.infer<typeof CardResponseSchema>;

// GET /api/v1/cards/:id/render - Stream card image (PNG/SVG)
export const RenderCardRequestSchema = z.object({
  params: z.object({ id: z.string().uuid() }),
  query: z.object({
    format: z.enum(["png", "svg"]).default("png"),
  }),
});

// POST /api/v1/cards/analyze - Analyze card for tags and metadata
export const AnalyzeCardRequestSchema = CardDataSchema;
export type AnalyzeCardRequest = z.infer<typeof AnalyzeCardRequestSchema>;

export const AnalyzeCardResponseSchema = z.object({
  tags: z.array(z.string()),
  synergyScore: z.number().min(0).max(10),
  tribalTags: z.array(z.string()),
  metadata: z.record(z.string()),
});
export type AnalyzeCardResponse = z.infer<typeof AnalyzeCardResponseSchema>;

// POST /api/v1/import/csv - Import cards from CSV (async job)
export const ImportCsvRequestSchema = z.object({
  file: z.instanceof(File), // Handled as multipart/form-data
  cardType: z
    .enum(["creature", "spell", "artifact", "incantation", "anthem"])
    .optional(),
  dryRun: z.boolean().default(false),
});

export const ImportCsvResponseSchema = z.object({
  jobId: z.string().uuid(),
  status: z.enum(["started", "processing", "completed", "failed"]),
});
export type ImportCsvResponse = z.infer<typeof ImportCsvResponseSchema>;

// GET /api/v1/import/:jobId/status - Check import job status
export const GetImportStatusRequestSchema = z.object({
  params: z.object({ jobId: z.string().uuid() }),
});

export const ImportStatusResponseSchema = z.object({
  jobId: z.string().uuid(),
  status: z.enum(["processing", "completed", "failed"]),
  importedCount: z.number(),
  totalCount: z.number(),
  errors: z.array(z.string()),
});
export type ImportStatusResponse = z.infer<typeof ImportStatusResponseSchema>;

// DELETE /api/v1/admin/cards - Clear all cards (admin endpoint)
export const ClearCardsResponseSchema = z.object({
  clearedCount: z.number(),
});

// GET /api/v1/health - Service health status
export const HealthResponseSchema = z.object({
  status: z.enum(["ok", "degraded", "error"]),
  services: z.record(
    z.object({
      status: z.enum(["ok", "degraded", "error"]),
      uptime: z.string().optional(),
    }),
  ),
});
export type HealthResponse = z.infer<typeof HealthResponseSchema>;
