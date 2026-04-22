/**
 * Copyright 2026 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { InspectionDataBuilder } from 'src/app/parser/core/builder';
import { IDataAssembler } from 'src/app/parser/core/interfaces';

/**
 * A dummy implementation of IDataAssembler for testing or placeholders.
 */
export class V6DummyAssembler<
  TProto = unknown,
> implements IDataAssembler<TProto> {
  constructor(private readonly name: string) {}

  /**
   * Ingests a decoded Protobuf chunk.
   */
  ingest(proto: TProto): void {
    console.log(`[${this.name}] Ingesting proto:`, proto);
  }

  /**
   * Integrates the ingested data into the InspectionDataBuilder.
   */
  assembleInto(builder: InspectionDataBuilder): void {
    console.log(`[${this.name}] Assembling into builder.`);
    // TODO: Mutate the builder when implemented
  }
}
