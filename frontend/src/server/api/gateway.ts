import { TRPCError } from "@trpc/server";
import { randomUUID } from "crypto";

import type {
  GenerateCardRequest,
  GenerateCardResponse,
  CardResponse,
  AnalyzeCardRequest,
  AnalyzeCardResponse,
  ImportCsvResponse,
  ImportStatusResponse,
  HealthResponse,
} from "../../../../api/contracts";

import {
  GenerateCardResponseSchema,
  CardResponseSchema,
  AnalyzeCardResponseSchema,
  ImportCsvResponseSchema,
  ImportStatusResponseSchema,
  ClearCardsResponseSchema,
  HealthResponseSchema,
} from "../../../../api/contracts";

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

const BASE_URL =
  process.env.API_GATEWAY_URL ?? "http://api-gateway:8080/api/v1";

const TIMEOUTS = {
  card: 30_000,
  import: 60_000,
  health: 5_000,
} as const;

const MAX_RETRIES = 3;
const RETRYABLE_STATUS_CODES = new Set([429, 503]);

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function trpcCodeFromStatus(
  status: number,
): ConstructorParameters<typeof TRPCError>[0]["code"] {
  if (status === 400) return "BAD_REQUEST";
  if (status === 401) return "UNAUTHORIZED";
  if (status === 403) return "FORBIDDEN";
  if (status === 404) return "NOT_FOUND";
  if (status === 408) return "TIMEOUT";
  if (status === 409) return "CONFLICT";
  if (status === 429) return "TOO_MANY_REQUESTS";
  if (status >= 500) return "INTERNAL_SERVER_ERROR";
  return "INTERNAL_SERVER_ERROR";
}

async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// ---------------------------------------------------------------------------
// Gateway Client
// ---------------------------------------------------------------------------

export class GatewayClient {
  private readonly baseUrl: string;

  constructor(baseUrl: string = BASE_URL) {
    this.baseUrl = baseUrl.replace(/\/+$/, "");
  }

  // -----------------------------------------------------------------------
  // Core fetch with retry + timeout
  // -----------------------------------------------------------------------

  private async request<T>(
    path: string,
    options: {
      method?: string;
      body?: unknown;
      formData?: FormData;
      timeout: number;
    },
  ): Promise<T> {
    const url = `${this.baseUrl}${path}`;
    const requestId = randomUUID();
    const { method = "GET", body, formData, timeout } = options;

    const headers: Record<string, string> = {
      "X-Request-ID": requestId,
    };

    let fetchBody: BodyInit | undefined;

    if (formData) {
      // Let the runtime set the multipart boundary via Content-Type
      fetchBody = formData;
    } else if (body !== undefined) {
      headers["Content-Type"] = "application/json";
      fetchBody = JSON.stringify(body);
    }

    let lastError: unknown;

    for (let attempt = 0; attempt <= MAX_RETRIES; attempt++) {
      const controller = new AbortController();
      const timer = setTimeout(() => controller.abort(), timeout);

      try {
        const res = await fetch(url, {
          method,
          headers,
          body: fetchBody,
          signal: controller.signal,
        });

        clearTimeout(timer);

        if (res.ok) {
          return (await res.json()) as T;
        }

        // Retry on 429 / 503 if we have attempts left
        if (RETRYABLE_STATUS_CODES.has(res.status) && attempt < MAX_RETRIES) {
          const backoff = Math.pow(2, attempt) * 500; // 500ms, 1s, 2s
          await sleep(backoff);
          lastError = new TRPCError({
            code: trpcCodeFromStatus(res.status),
            message: `Gateway ${method} ${path} returned ${res.status}`,
          });
          continue;
        }

        // Non-retryable or retries exhausted
        let detail: string | undefined;
        try {
          const errBody = (await res.json()) as { error?: string; message?: string };
          detail = errBody.error ?? errBody.message;
        } catch {
          // body wasn't JSON – ignore
        }

        throw new TRPCError({
          code: trpcCodeFromStatus(res.status),
          message: detail ?? `Gateway ${method} ${path} returned ${res.status}`,
        });
      } catch (err) {
        clearTimeout(timer);

        if (err instanceof TRPCError) {
          // If it's already a TRPCError from the block above, rethrow unless retryable
          if (attempt >= MAX_RETRIES) throw err;
          lastError = err;
          continue;
        }

        if ((err as Error).name === "AbortError") {
          throw new TRPCError({
            code: "TIMEOUT",
            message: `Gateway ${method} ${path} timed out after ${timeout}ms (request-id: ${requestId})`,
          });
        }

        // Network / DNS errors
        if (attempt < MAX_RETRIES) {
          const backoff = Math.pow(2, attempt) * 500;
          await sleep(backoff);
          lastError = err;
          continue;
        }

        throw new TRPCError({
          code: "INTERNAL_SERVER_ERROR",
          message: `Gateway ${method} ${path} failed: ${(err as Error).message}`,
          cause: err,
        });
      }
    }

    // Should be unreachable, but just in case
    throw (
      lastError ??
      new TRPCError({
        code: "INTERNAL_SERVER_ERROR",
        message: "Unexpected retry exhaustion",
      })
    );
  }

  // -----------------------------------------------------------------------
  // Endpoints
  // -----------------------------------------------------------------------

  /** GET /health */
  async health(): Promise<HealthResponse> {
    const raw = await this.request<unknown>("/health", {
      timeout: TIMEOUTS.health,
    });
    return HealthResponseSchema.parse(raw);
  }

  /** POST /cards/generate */
  async generateCard(data: GenerateCardRequest): Promise<GenerateCardResponse> {
    const raw = await this.request<unknown>("/cards/generate", {
      method: "POST",
      body: data,
      timeout: TIMEOUTS.card,
    });
    return GenerateCardResponseSchema.parse(raw);
  }

  /** GET /cards/:id */
  async getCard(id: string): Promise<CardResponse> {
    const raw = await this.request<unknown>(`/cards/${encodeURIComponent(id)}`, {
      timeout: TIMEOUTS.card,
    });
    return CardResponseSchema.parse(raw);
  }

  /** POST /cards/analyze */
  async analyzeCard(data: AnalyzeCardRequest): Promise<AnalyzeCardResponse> {
    const raw = await this.request<unknown>("/cards/analyze", {
      method: "POST",
      body: data,
      timeout: TIMEOUTS.card,
    });
    return AnalyzeCardResponseSchema.parse(raw);
  }

  /** POST /import/csv */
  async importCsv(
    file: File,
    cardType?: string,
    dryRun?: boolean,
  ): Promise<ImportCsvResponse> {
    const formData = new FormData();
    formData.append("file", file);
    if (cardType) formData.append("cardType", cardType);
    if (dryRun !== undefined) formData.append("dryRun", String(dryRun));

    const raw = await this.request<unknown>("/import/csv", {
      method: "POST",
      formData,
      timeout: TIMEOUTS.import,
    });
    return ImportCsvResponseSchema.parse(raw);
  }

  /** GET /import/:jobId/status */
  async getImportStatus(jobId: string): Promise<ImportStatusResponse> {
    const raw = await this.request<unknown>(
      `/import/${encodeURIComponent(jobId)}/status`,
      { timeout: TIMEOUTS.import },
    );
    return ImportStatusResponseSchema.parse(raw);
  }

  /** DELETE /admin/cards */
  async clearCards(): Promise<{ clearedCount: number }> {
    const raw = await this.request<unknown>("/admin/cards", {
      method: "DELETE",
      timeout: TIMEOUTS.card,
    });
    return ClearCardsResponseSchema.parse(raw);
  }
}

// ---------------------------------------------------------------------------
// Singleton
// ---------------------------------------------------------------------------

export const gateway = new GatewayClient();
