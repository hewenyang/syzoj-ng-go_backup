// package: syzoj.model
// file: syzoj.model.proto

import * as jspb from "google-protobuf";

export class UserAuth extends jspb.Message {
  hasPasswordHash(): boolean;
  clearPasswordHash(): void;
  getPasswordHash(): Uint8Array | string;
  getPasswordHash_asU8(): Uint8Array;
  getPasswordHash_asB64(): string;
  setPasswordHash(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAuth.AsObject;
  static toObject(includeInstance: boolean, msg: UserAuth): UserAuth.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserAuth;
  static deserializeBinaryFromReader(message: UserAuth, reader: jspb.BinaryReader): UserAuth;
}

export namespace UserAuth {
  export type AsObject = {
    passwordHash: Uint8Array | string,
  }
}

export class DeviceInfo extends jspb.Message {
  hasToken(): boolean;
  clearToken(): void;
  getToken(): string | undefined;
  setToken(value: string): void;

  hasUserAgent(): boolean;
  clearUserAgent(): void;
  getUserAgent(): string | undefined;
  setUserAgent(value: string): void;

  hasRemoteAddr(): boolean;
  clearRemoteAddr(): void;
  getRemoteAddr(): string | undefined;
  setRemoteAddr(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeviceInfo.AsObject;
  static toObject(includeInstance: boolean, msg: DeviceInfo): DeviceInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeviceInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeviceInfo;
  static deserializeBinaryFromReader(message: DeviceInfo, reader: jspb.BinaryReader): DeviceInfo;
}

export namespace DeviceInfo {
  export type AsObject = {
    token?: string,
    userAgent?: string,
    remoteAddr?: string,
  }
}

export class Problem extends jspb.Message {
  hasTitle(): boolean;
  clearTitle(): void;
  getTitle(): string | undefined;
  setTitle(value: string): void;

  hasStatement(): boolean;
  clearStatement(): void;
  getStatement(): ProblemStatement | undefined;
  setStatement(value?: ProblemStatement): void;

  clearSourceList(): void;
  getSourceList(): Array<ProblemSource>;
  setSourceList(value: Array<ProblemSource>): void;
  addSource(value?: ProblemSource, index?: number): ProblemSource;

  hasJudge(): boolean;
  clearJudge(): void;
  getJudge(): ProblemJudge | undefined;
  setJudge(value?: ProblemJudge): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Problem.AsObject;
  static toObject(includeInstance: boolean, msg: Problem): Problem.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Problem, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Problem;
  static deserializeBinaryFromReader(message: Problem, reader: jspb.BinaryReader): Problem;
}

export namespace Problem {
  export type AsObject = {
    title?: string,
    statement?: ProblemStatement.AsObject,
    sourceList: Array<ProblemSource.AsObject>,
    judge?: ProblemJudge.AsObject,
  }
}

export class ProblemJudge extends jspb.Message {
  hasTraditional(): boolean;
  clearTraditional(): void;
  getTraditional(): TraditionalJudgeData | undefined;
  setTraditional(value?: TraditionalJudgeData): void;

  getJudgeCase(): ProblemJudge.JudgeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProblemJudge.AsObject;
  static toObject(includeInstance: boolean, msg: ProblemJudge): ProblemJudge.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProblemJudge, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProblemJudge;
  static deserializeBinaryFromReader(message: ProblemJudge, reader: jspb.BinaryReader): ProblemJudge;
}

export namespace ProblemJudge {
  export type AsObject = {
    traditional?: TraditionalJudgeData.AsObject,
  }

  export enum JudgeCase {
    JUDGE_NOT_SET = 0,
    TRADITIONAL = 1,
  }
}

export class ProblemSource extends jspb.Message {
  hasUrl(): boolean;
  clearUrl(): void;
  getUrl(): string | undefined;
  setUrl(value: string): void;

  hasSiteName(): boolean;
  clearSiteName(): void;
  getSiteName(): string | undefined;
  setSiteName(value: string): void;

  hasSiteProblemTitle(): boolean;
  clearSiteProblemTitle(): void;
  getSiteProblemTitle(): string | undefined;
  setSiteProblemTitle(value: string): void;

  hasContestName(): boolean;
  clearContestName(): void;
  getContestName(): string | undefined;
  setContestName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProblemSource.AsObject;
  static toObject(includeInstance: boolean, msg: ProblemSource): ProblemSource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProblemSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProblemSource;
  static deserializeBinaryFromReader(message: ProblemSource, reader: jspb.BinaryReader): ProblemSource;
}

export namespace ProblemSource {
  export type AsObject = {
    url?: string,
    siteName?: string,
    siteProblemTitle?: string,
    contestName?: string,
  }
}

export class ProblemStatement extends jspb.Message {
  hasMarkdown(): boolean;
  clearMarkdown(): void;
  getMarkdown(): MarkdownLatexDocument | undefined;
  setMarkdown(value?: MarkdownLatexDocument): void;

  getStatementCase(): ProblemStatement.StatementCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProblemStatement.AsObject;
  static toObject(includeInstance: boolean, msg: ProblemStatement): ProblemStatement.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProblemStatement, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProblemStatement;
  static deserializeBinaryFromReader(message: ProblemStatement, reader: jspb.BinaryReader): ProblemStatement;
}

export namespace ProblemStatement {
  export type AsObject = {
    markdown?: MarkdownLatexDocument.AsObject,
  }

  export enum StatementCase {
    STATEMENT_NOT_SET = 0,
    MARKDOWN = 16,
  }
}

export class MarkdownLatexDocument extends jspb.Message {
  hasText(): boolean;
  clearText(): void;
  getText(): string | undefined;
  setText(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MarkdownLatexDocument.AsObject;
  static toObject(includeInstance: boolean, msg: MarkdownLatexDocument): MarkdownLatexDocument.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MarkdownLatexDocument, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MarkdownLatexDocument;
  static deserializeBinaryFromReader(message: MarkdownLatexDocument, reader: jspb.BinaryReader): MarkdownLatexDocument;
}

export namespace MarkdownLatexDocument {
  export type AsObject = {
    text?: string,
  }
}

export class TraditionalJudgeData extends jspb.Message {
  hasTimeLimit(): boolean;
  clearTimeLimit(): void;
  getTimeLimit(): number | undefined;
  setTimeLimit(value: number): void;

  hasMemoryLimit(): boolean;
  clearMemoryLimit(): void;
  getMemoryLimit(): number | undefined;
  setMemoryLimit(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TraditionalJudgeData.AsObject;
  static toObject(includeInstance: boolean, msg: TraditionalJudgeData): TraditionalJudgeData.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TraditionalJudgeData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TraditionalJudgeData;
  static deserializeBinaryFromReader(message: TraditionalJudgeData, reader: jspb.BinaryReader): TraditionalJudgeData;
}

export namespace TraditionalJudgeData {
  export type AsObject = {
    timeLimit?: number,
    memoryLimit?: number,
  }
}

export class TraditionalJudgeCode extends jspb.Message {
  hasLanguage(): boolean;
  clearLanguage(): void;
  getLanguage(): string | undefined;
  setLanguage(value: string): void;

  hasCode(): boolean;
  clearCode(): void;
  getCode(): string | undefined;
  setCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TraditionalJudgeCode.AsObject;
  static toObject(includeInstance: boolean, msg: TraditionalJudgeCode): TraditionalJudgeCode.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TraditionalJudgeCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TraditionalJudgeCode;
  static deserializeBinaryFromReader(message: TraditionalJudgeCode, reader: jspb.BinaryReader): TraditionalJudgeCode;
}

export namespace TraditionalJudgeCode {
  export type AsObject = {
    language?: string,
    code?: string,
  }
}

