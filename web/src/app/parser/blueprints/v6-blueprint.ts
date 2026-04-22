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

import { fromBinary } from '@bufbuild/protobuf';
import { ParserBlueprint } from 'src/app/parser/core/interfaces';
import { V6DummyAssembler } from 'src/app/parser/assemblers/v6/dummy-assemblers';

import { MetadataChunkSchema } from 'src/app/generated/khifile/v6/metadata_pb';
import { InterningPoolChunkSchema } from 'src/app/generated/khifile/v6/intern_pool_pb';
import { LogChunkSchema } from 'src/app/generated/khifile/v6/log_pb';
import { TimelineStyleChunkSchema } from 'src/app/generated/khifile/v6/style_pb';
import { TimelineChunkSchema } from 'src/app/generated/khifile/v6/timeline_pb';

export enum ChunkType {
  METADATA = 1,
  INTERN_POOL = 2,
  LOG = 3,
  TIMELINE_STYLE = 4,
  TIMELINE = 5,
}

/**
 * Blueprint for V6 KHI File parsing.
 * Maps chunk type IDs to their respective decoding and assembly logic.
 */
export const V6_BLUEPRINT: ParserBlueprint = new Map([
  [
    ChunkType.METADATA,
    {
      typeId: ChunkType.METADATA,
      decode: (bytes) => fromBinary(MetadataChunkSchema, bytes),
      createAssembler: () => new V6DummyAssembler('MetadataAssembler'),
      priority: 0,
    },
  ],
  [
    ChunkType.INTERN_POOL,
    {
      typeId: ChunkType.INTERN_POOL,
      decode: (bytes) => fromBinary(InterningPoolChunkSchema, bytes),
      createAssembler: () => new V6DummyAssembler('InterningPoolAssembler'),
      priority: 10,
    },
  ],
  [
    ChunkType.TIMELINE_STYLE,
    {
      typeId: ChunkType.TIMELINE_STYLE,
      decode: (bytes) => fromBinary(TimelineStyleChunkSchema, bytes),
      createAssembler: () => new V6DummyAssembler('TimelineStyleAssembler'),
      priority: 20,
    },
  ],
  [
    ChunkType.LOG,
    {
      typeId: ChunkType.LOG,
      decode: (bytes) => fromBinary(LogChunkSchema, bytes),
      createAssembler: () => new V6DummyAssembler('LogAssembler'),
      priority: 100, // Depends on string pool and style
    },
  ],
  [
    ChunkType.TIMELINE,
    {
      typeId: ChunkType.TIMELINE,
      decode: (bytes) => fromBinary(TimelineChunkSchema, bytes),
      createAssembler: () => new V6DummyAssembler('TimelineAssembler'),
      priority: 200, // Depends on logs, string pool, and style
    },
  ],
]);
