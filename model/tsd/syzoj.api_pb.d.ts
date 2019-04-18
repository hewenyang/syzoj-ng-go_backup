// package: syzoj.api
// file: syzoj.api.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as syzoj_model_pb from "./syzoj.model_pb";

export class Response extends jspb.Message {
  clearMutationsList(): void;
  getMutationsList(): Array<Mutation>;
  setMutationsList(value: Array<Mutation>): void;
  addMutations(value?: Mutation, index?: number): Mutation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Response.AsObject;
  static toObject(includeInstance: boolean, msg: Response): Response.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Response;
  static deserializeBinaryFromReader(message: Response, reader: jspb.BinaryReader): Response;
}

export namespace Response {
  export type AsObject = {
    mutationsList: Array<Mutation.AsObject>,
  }
}

export class Mutation extends jspb.Message {
  hasPath(): boolean;
  clearPath(): void;
  getPath(): string | undefined;
  setPath(value: string): void;

  hasMethod(): boolean;
  clearMethod(): void;
  getMethod(): string | undefined;
  setMethod(value: string): void;

  hasValue(): boolean;
  clearValue(): void;
  getValue(): google_protobuf_any_pb.Any | undefined;
  setValue(value?: google_protobuf_any_pb.Any): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Mutation.AsObject;
  static toObject(includeInstance: boolean, msg: Mutation): Mutation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Mutation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Mutation;
  static deserializeBinaryFromReader(message: Mutation, reader: jspb.BinaryReader): Mutation;
}

export namespace Mutation {
  export type AsObject = {
    path?: string,
    method?: string,
    value?: google_protobuf_any_pb.Any.AsObject,
  }
}

export class Error extends jspb.Message {
  hasError(): boolean;
  clearError(): void;
  getError(): string | undefined;
  setError(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Error.AsObject;
  static toObject(includeInstance: boolean, msg: Error): Error.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Error, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Error;
  static deserializeBinaryFromReader(message: Error, reader: jspb.BinaryReader): Error;
}

export namespace Error {
  export type AsObject = {
    error?: string,
  }
}

export class Path extends jspb.Message {
  hasPath(): boolean;
  clearPath(): void;
  getPath(): string | undefined;
  setPath(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Path.AsObject;
  static toObject(includeInstance: boolean, msg: Path): Path.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Path, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Path;
  static deserializeBinaryFromReader(message: Path, reader: jspb.BinaryReader): Path;
}

export namespace Path {
  export type AsObject = {
    path?: string,
  }
}

export class NotFoundPage extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotFoundPage.AsObject;
  static toObject(includeInstance: boolean, msg: NotFoundPage): NotFoundPage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NotFoundPage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotFoundPage;
  static deserializeBinaryFromReader(message: NotFoundPage, reader: jspb.BinaryReader): NotFoundPage;
}

export namespace NotFoundPage {
  export type AsObject = {
  }
}

export class IndexPage extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IndexPage.AsObject;
  static toObject(includeInstance: boolean, msg: IndexPage): IndexPage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: IndexPage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IndexPage;
  static deserializeBinaryFromReader(message: IndexPage, reader: jspb.BinaryReader): IndexPage;
}

export namespace IndexPage {
  export type AsObject = {
  }
}

export class LoginPage extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginPage.AsObject;
  static toObject(includeInstance: boolean, msg: LoginPage): LoginPage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LoginPage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginPage;
  static deserializeBinaryFromReader(message: LoginPage, reader: jspb.BinaryReader): LoginPage;
}

export namespace LoginPage {
  export type AsObject = {
  }

  export class LoginRequest extends jspb.Message {
    hasUserName(): boolean;
    clearUserName(): void;
    getUserName(): string | undefined;
    setUserName(value: string): void;

    hasPassword(): boolean;
    clearPassword(): void;
    getPassword(): string | undefined;
    setPassword(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LoginRequest.AsObject;
    static toObject(includeInstance: boolean, msg: LoginRequest): LoginRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LoginRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LoginRequest;
    static deserializeBinaryFromReader(message: LoginRequest, reader: jspb.BinaryReader): LoginRequest;
  }

  export namespace LoginRequest {
    export type AsObject = {
      userName?: string,
      password?: string,
    }
  }
}

export class RegisterPage extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterPage.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterPage): RegisterPage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RegisterPage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterPage;
  static deserializeBinaryFromReader(message: RegisterPage, reader: jspb.BinaryReader): RegisterPage;
}

export namespace RegisterPage {
  export type AsObject = {
  }

  export class RegisterRequest extends jspb.Message {
    hasUserName(): boolean;
    clearUserName(): void;
    getUserName(): string | undefined;
    setUserName(value: string): void;

    hasPassword(): boolean;
    clearPassword(): void;
    getPassword(): string | undefined;
    setPassword(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RegisterRequest.AsObject;
    static toObject(includeInstance: boolean, msg: RegisterRequest): RegisterRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RegisterRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RegisterRequest;
    static deserializeBinaryFromReader(message: RegisterRequest, reader: jspb.BinaryReader): RegisterRequest;
  }

  export namespace RegisterRequest {
    export type AsObject = {
      userName?: string,
      password?: string,
    }
  }
}

export class ProblemCreatePage extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProblemCreatePage.AsObject;
  static toObject(includeInstance: boolean, msg: ProblemCreatePage): ProblemCreatePage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProblemCreatePage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProblemCreatePage;
  static deserializeBinaryFromReader(message: ProblemCreatePage, reader: jspb.BinaryReader): ProblemCreatePage;
}

export namespace ProblemCreatePage {
  export type AsObject = {
  }

  export class CreateRequest extends jspb.Message {
    hasProblemTitle(): boolean;
    clearProblemTitle(): void;
    getProblemTitle(): string | undefined;
    setProblemTitle(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CreateRequest.AsObject;
    static toObject(includeInstance: boolean, msg: CreateRequest): CreateRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CreateRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CreateRequest;
    static deserializeBinaryFromReader(message: CreateRequest, reader: jspb.BinaryReader): CreateRequest;
  }

  export namespace CreateRequest {
    export type AsObject = {
      problemTitle?: string,
    }
  }
}

export class ProblemViewPage extends jspb.Message {
  hasProblemTitle(): boolean;
  clearProblemTitle(): void;
  getProblemTitle(): string | undefined;
  setProblemTitle(value: string): void;

  hasProblemStatement(): boolean;
  clearProblemStatement(): void;
  getProblemStatement(): syzoj_model_pb.ProblemStatement | undefined;
  setProblemStatement(value?: syzoj_model_pb.ProblemStatement): void;

  clearProblemSourceList(): void;
  getProblemSourceList(): Array<syzoj_model_pb.ProblemSource>;
  setProblemSourceList(value: Array<syzoj_model_pb.ProblemSource>): void;
  addProblemSource(value?: syzoj_model_pb.ProblemSource, index?: number): syzoj_model_pb.ProblemSource;

  hasProblemJudge(): boolean;
  clearProblemJudge(): void;
  getProblemJudge(): syzoj_model_pb.ProblemJudge | undefined;
  setProblemJudge(value?: syzoj_model_pb.ProblemJudge): void;

  clearProblemEntryList(): void;
  getProblemEntryList(): Array<ProblemViewPage.ProblemEntry>;
  setProblemEntryList(value: Array<ProblemViewPage.ProblemEntry>): void;
  addProblemEntry(value?: ProblemViewPage.ProblemEntry, index?: number): ProblemViewPage.ProblemEntry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProblemViewPage.AsObject;
  static toObject(includeInstance: boolean, msg: ProblemViewPage): ProblemViewPage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProblemViewPage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProblemViewPage;
  static deserializeBinaryFromReader(message: ProblemViewPage, reader: jspb.BinaryReader): ProblemViewPage;
}

export namespace ProblemViewPage {
  export type AsObject = {
    problemTitle?: string,
    problemStatement?: syzoj_model_pb.ProblemStatement.AsObject,
    problemSourceList: Array<syzoj_model_pb.ProblemSource.AsObject>,
    problemJudge?: syzoj_model_pb.ProblemJudge.AsObject,
    problemEntryList: Array<ProblemViewPage.ProblemEntry.AsObject>,
  }

  export class ProblemEntry extends jspb.Message {
    hasId(): boolean;
    clearId(): void;
    getId(): string | undefined;
    setId(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ProblemEntry.AsObject;
    static toObject(includeInstance: boolean, msg: ProblemEntry): ProblemEntry.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ProblemEntry, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ProblemEntry;
    static deserializeBinaryFromReader(message: ProblemEntry, reader: jspb.BinaryReader): ProblemEntry;
  }

  export namespace ProblemEntry {
    export type AsObject = {
      id?: string,
    }
  }

  export class AddStatementRequest extends jspb.Message {
    hasStatement(): boolean;
    clearStatement(): void;
    getStatement(): syzoj_model_pb.ProblemStatement | undefined;
    setStatement(value?: syzoj_model_pb.ProblemStatement): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AddStatementRequest.AsObject;
    static toObject(includeInstance: boolean, msg: AddStatementRequest): AddStatementRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AddStatementRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AddStatementRequest;
    static deserializeBinaryFromReader(message: AddStatementRequest, reader: jspb.BinaryReader): AddStatementRequest;
  }

  export namespace AddStatementRequest {
    export type AsObject = {
      statement?: syzoj_model_pb.ProblemStatement.AsObject,
    }
  }

  export class AddSourceRequest extends jspb.Message {
    hasSource(): boolean;
    clearSource(): void;
    getSource(): syzoj_model_pb.ProblemSource | undefined;
    setSource(value?: syzoj_model_pb.ProblemSource): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AddSourceRequest.AsObject;
    static toObject(includeInstance: boolean, msg: AddSourceRequest): AddSourceRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AddSourceRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AddSourceRequest;
    static deserializeBinaryFromReader(message: AddSourceRequest, reader: jspb.BinaryReader): AddSourceRequest;
  }

  export namespace AddSourceRequest {
    export type AsObject = {
      source?: syzoj_model_pb.ProblemSource.AsObject,
    }
  }

  export class SetPublicRequest extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SetPublicRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SetPublicRequest): SetPublicRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SetPublicRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SetPublicRequest;
    static deserializeBinaryFromReader(message: SetPublicRequest, reader: jspb.BinaryReader): SetPublicRequest;
  }

  export namespace SetPublicRequest {
    export type AsObject = {
    }
  }

  export class AddJudgeTraditionalRequest extends jspb.Message {
    hasData(): boolean;
    clearData(): void;
    getData(): syzoj_model_pb.TraditionalJudgeData | undefined;
    setData(value?: syzoj_model_pb.TraditionalJudgeData): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AddJudgeTraditionalRequest.AsObject;
    static toObject(includeInstance: boolean, msg: AddJudgeTraditionalRequest): AddJudgeTraditionalRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AddJudgeTraditionalRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AddJudgeTraditionalRequest;
    static deserializeBinaryFromReader(message: AddJudgeTraditionalRequest, reader: jspb.BinaryReader): AddJudgeTraditionalRequest;
  }

  export namespace AddJudgeTraditionalRequest {
    export type AsObject = {
      data?: syzoj_model_pb.TraditionalJudgeData.AsObject,
    }
  }

  export class SubmitJudgeTraditionalRequest extends jspb.Message {
    hasCode(): boolean;
    clearCode(): void;
    getCode(): syzoj_model_pb.TraditionalJudgeCode | undefined;
    setCode(value?: syzoj_model_pb.TraditionalJudgeCode): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SubmitJudgeTraditionalRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SubmitJudgeTraditionalRequest): SubmitJudgeTraditionalRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SubmitJudgeTraditionalRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SubmitJudgeTraditionalRequest;
    static deserializeBinaryFromReader(message: SubmitJudgeTraditionalRequest, reader: jspb.BinaryReader): SubmitJudgeTraditionalRequest;
  }

  export namespace SubmitJudgeTraditionalRequest {
    export type AsObject = {
      code?: syzoj_model_pb.TraditionalJudgeCode.AsObject,
    }
  }

  export class SubmitJudgeTraditionalResponse extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SubmitJudgeTraditionalResponse.AsObject;
    static toObject(includeInstance: boolean, msg: SubmitJudgeTraditionalResponse): SubmitJudgeTraditionalResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SubmitJudgeTraditionalResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SubmitJudgeTraditionalResponse;
    static deserializeBinaryFromReader(message: SubmitJudgeTraditionalResponse, reader: jspb.BinaryReader): SubmitJudgeTraditionalResponse;
  }

  export namespace SubmitJudgeTraditionalResponse {
    export type AsObject = {
    }
  }
}

export class ProblemsPage extends jspb.Message {
  clearProblemEntryList(): void;
  getProblemEntryList(): Array<ProblemsPage.ProblemEntry>;
  setProblemEntryList(value: Array<ProblemsPage.ProblemEntry>): void;
  addProblemEntry(value?: ProblemsPage.ProblemEntry, index?: number): ProblemsPage.ProblemEntry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProblemsPage.AsObject;
  static toObject(includeInstance: boolean, msg: ProblemsPage): ProblemsPage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProblemsPage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProblemsPage;
  static deserializeBinaryFromReader(message: ProblemsPage, reader: jspb.BinaryReader): ProblemsPage;
}

export namespace ProblemsPage {
  export type AsObject = {
    problemEntryList: Array<ProblemsPage.ProblemEntry.AsObject>,
  }

  export class ProblemEntry extends jspb.Message {
    hasId(): boolean;
    clearId(): void;
    getId(): string | undefined;
    setId(value: string): void;

    hasProblemId(): boolean;
    clearProblemId(): void;
    getProblemId(): string | undefined;
    setProblemId(value: string): void;

    hasProblemTitle(): boolean;
    clearProblemTitle(): void;
    getProblemTitle(): string | undefined;
    setProblemTitle(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ProblemEntry.AsObject;
    static toObject(includeInstance: boolean, msg: ProblemEntry): ProblemEntry.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ProblemEntry, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ProblemEntry;
    static deserializeBinaryFromReader(message: ProblemEntry, reader: jspb.BinaryReader): ProblemEntry;
  }

  export namespace ProblemEntry {
    export type AsObject = {
      id?: string,
      problemId?: string,
      problemTitle?: string,
    }
  }

  export class AddProblemRequest extends jspb.Message {
    hasProblemId(): boolean;
    clearProblemId(): void;
    getProblemId(): string | undefined;
    setProblemId(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AddProblemRequest.AsObject;
    static toObject(includeInstance: boolean, msg: AddProblemRequest): AddProblemRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AddProblemRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AddProblemRequest;
    static deserializeBinaryFromReader(message: AddProblemRequest, reader: jspb.BinaryReader): AddProblemRequest;
  }

  export namespace AddProblemRequest {
    export type AsObject = {
      problemId?: string,
    }
  }
}

export class DebugAddJudgerRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DebugAddJudgerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DebugAddJudgerRequest): DebugAddJudgerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DebugAddJudgerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DebugAddJudgerRequest;
  static deserializeBinaryFromReader(message: DebugAddJudgerRequest, reader: jspb.BinaryReader): DebugAddJudgerRequest;
}

export namespace DebugAddJudgerRequest {
  export type AsObject = {
  }
}

export class DebugAddJudgerResponse extends jspb.Message {
  hasJudgerId(): boolean;
  clearJudgerId(): void;
  getJudgerId(): string | undefined;
  setJudgerId(value: string): void;

  hasJudgerToken(): boolean;
  clearJudgerToken(): void;
  getJudgerToken(): string | undefined;
  setJudgerToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DebugAddJudgerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DebugAddJudgerResponse): DebugAddJudgerResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DebugAddJudgerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DebugAddJudgerResponse;
  static deserializeBinaryFromReader(message: DebugAddJudgerResponse, reader: jspb.BinaryReader): DebugAddJudgerResponse;
}

export namespace DebugAddJudgerResponse {
  export type AsObject = {
    judgerId?: string,
    judgerToken?: string,
  }
}

