// Copyright (C) 2019-2020 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License.

#include "db/insert/SSVectorSource.h"

#include <utility>
#include <vector>

#include "db/engine/EngineFactory.h"
#include "db/engine/ExecutionEngine.h"
#include "metrics/Metrics.h"
#include "utils/Log.h"
#include "utils/TimeRecorder.h"

namespace milvus {
namespace engine {

SSVectorSource::SSVectorSource(const DataChunkPtr& chunk) : chunk_(chunk) {
}

Status
SSVectorSource::Add(const segment::SSSegmentWriterPtr& segment_writer_ptr, const int64_t& num_entities_to_add,
                    int64_t& num_entities_added) {
    // TODO: n = vectors_.vector_count_;???
    int64_t n = chunk_->count_;
    num_entities_added = current_num_added_ + num_entities_to_add <= n ? num_entities_to_add : n - current_num_added_;

    auto status = segment_writer_ptr->AddChunk(chunk_, current_num_added_, num_entities_added);
    if (!status.ok()) {
        return status;
    }

    current_num_added_ += num_entities_added;
    return status;
}

bool
SSVectorSource::AllAdded() {
    return (current_num_added_ >= chunk_->count_);
}

}  // namespace engine
}  // namespace milvus
