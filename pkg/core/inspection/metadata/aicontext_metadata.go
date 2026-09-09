// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inspectionmetadata

import (
	"cmp"
	"slices"
	"strings"
	"sync"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	pb "github.com/GoogleCloudPlatform/khi/pkg/generated/khifile/v6"
	"google.golang.org/protobuf/proto"
)

// DefaultSectionPriority is the default sorting priority for an AI context section.
const DefaultSectionPriority = 1000

// aiContextSectionState holds the in-memory state of a single context section.
type aiContextSectionState struct {
	Title           string
	Priority        int32
	Properties      map[string]string
	SetProperties   map[string]map[string]struct{}
	SummaryMarkdown string
}

// AIContextSection represents an immutable read-only snapshot of a context section.
type AIContextSection struct {
	Title           string
	Priority        int32
	Properties      map[string]string
	SetProperties   map[string][]string
	SummaryMarkdown string
}

// AIContextMetadata collects contextual intelligence and summaries recorded by parsers.
type AIContextMetadata struct {
	sections map[string]*aiContextSectionState
	mu       sync.RWMutex
}

var _ Metadata = (*AIContextMetadata)(nil)

// NewAIContextMetadata creates a new instance of AIContextMetadata.
func NewAIContextMetadata() *AIContextMetadata {
	return &AIContextMetadata{
		sections: make(map[string]*aiContextSectionState),
	}
}

// Labels implements Metadata.
func (*AIContextMetadata) Labels() *typedmap.ReadonlyTypedMap {
	return NewLabelSet(IncludeInRunResult(), IncludeInTaskList(), IncludeInResultBinary())
}

// ToSerializable implements Metadata.
func (m *AIContextMetadata) ToSerializable() interface{} {
	snapshots := m.Sections()
	pbSections := make([]*pb.AIContextSection, len(snapshots))

	for i, snap := range snapshots {
		pbSetProps := make(map[string]*pb.StringList, len(snap.SetProperties))
		for k, vals := range snap.SetProperties {
			pbSetProps[k] = &pb.StringList{
				Values: vals,
			}
		}

		pbSections[i] = &pb.AIContextSection{
			Title:           proto.String(snap.Title),
			Priority:        proto.Int32(snap.Priority),
			Properties:      snap.Properties,
			SetProperties:   pbSetProps,
			SummaryMarkdown: proto.String(snap.SummaryMarkdown),
		}
	}

	return &pb.AIContextMetadata{
		Sections: pbSections,
	}
}

// SetPriority sets the display priority for a specific section.
func (m *AIContextMetadata) SetPriority(sectionTitle string, priority int32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sec := m.getOrCreateSectionLocked(sectionTitle)
	sec.Priority = priority
}

// SetProperty records or updates a scalar key-value property under the given section.
func (m *AIContextMetadata) SetProperty(sectionTitle, key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sec := m.getOrCreateSectionLocked(sectionTitle)
	sec.Properties[key] = value
}

// AddToSetProperty adds a string element into a set property under the given section.
func (m *AIContextMetadata) AddToSetProperty(sectionTitle, key, value string) {
	if value == "" {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	sec := m.getOrCreateSectionLocked(sectionTitle)
	set, exists := sec.SetProperties[key]
	if !exists {
		set = make(map[string]struct{})
		sec.SetProperties[key] = set
	}
	set[value] = struct{}{}
}

// AppendSummaryMarkdown appends markdown text separated by double newlines under the given section.
func (m *AIContextMetadata) AppendSummaryMarkdown(sectionTitle, markdown string) {
	trimmed := strings.Trim(markdown, "\n")
	if trimmed == "" {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	sec := m.getOrCreateSectionLocked(sectionTitle)
	if sec.SummaryMarkdown == "" {
		sec.SummaryMarkdown = trimmed
	} else {
		sec.SummaryMarkdown += "\n\n" + trimmed
	}
}

// Sections returns a sorted slice of immutable section snapshots.
func (m *AIContextMetadata) Sections() []*AIContextSection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snapshots := make([]*AIContextSection, 0, len(m.sections))
	for _, sec := range m.sections {
		propsCopy := make(map[string]string, len(sec.Properties))
		for k, v := range sec.Properties {
			propsCopy[k] = v
		}

		setPropsCopy := make(map[string][]string, len(sec.SetProperties))
		for k, set := range sec.SetProperties {
			vals := make([]string, 0, len(set))
			for val := range set {
				vals = append(vals, val)
			}
			slices.Sort(vals)
			setPropsCopy[k] = vals
		}

		snapshots = append(snapshots, &AIContextSection{
			Title:           sec.Title,
			Priority:        sec.Priority,
			Properties:      propsCopy,
			SetProperties:   setPropsCopy,
			SummaryMarkdown: sec.SummaryMarkdown,
		})
	}

	slices.SortFunc(snapshots, func(a, b *AIContextSection) int {
		if a.Priority != b.Priority {
			return cmp.Compare(a.Priority, b.Priority)
		}
		return cmp.Compare(a.Title, b.Title)
	})

	return snapshots
}

// getOrCreateSectionLocked retrieves an existing section or creates a new one.
// The caller must hold m.mu write lock.
func (m *AIContextMetadata) getOrCreateSectionLocked(sectionTitle string) *aiContextSectionState {
	sec, exists := m.sections[sectionTitle]
	if !exists {
		sec = &aiContextSectionState{
			Title:         sectionTitle,
			Priority:      DefaultSectionPriority,
			Properties:    make(map[string]string),
			SetProperties: make(map[string]map[string]struct{}),
		}
		m.sections[sectionTitle] = sec
	}
	return sec
}
