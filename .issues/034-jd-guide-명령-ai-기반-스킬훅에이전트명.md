---
number: 34
title: 'feat: jd guide 명령 - AI 기반 스킬/훅/에이전트/명령 사용법 안내'
state: done
labels:
  - feature
assignees: []
created_at: '2026-01-21T14:56:51Z'
updated_at: '2026-01-21T15:44:27Z'
closed_at: '2026-01-21T15:44:27Z'
---

## Background

- Claude CLI의 스킬, 훅, 에이전트, 명령을 설치해도 **언제 호출되고 어떻게 활용해야 하는지** 파악하기 어려움
- 기존 `show` 명령은 원본 마크다운만 보여주고, 실제 활용법 해석은 사용자 몫
- 다른 사람이 만든 스킬의 의도와 활용 시나리오를 빠르게 이해하고 싶은 니즈

## Problem

스킬/훅/에이전트/명령의 **사용법, 호출 시점, 활용 방법, 개선점**을 쉽게 알 수 없음

## Goal

**`jd guide <type> <id>` 명령으로 AI가 분석해서 사용법을 안내**

1. **기본 모드**: 1회성 설명 (박스 UI로 로딩 중 표시 → 결과 출력)
2. **대화형 모드** (`-i`): AI가 사용자 상황 파악 질문 후 맞춤형 안내 제공

### 지원 타입

- `jd guide skills <id>`
- `jd guide hooks <name>`
- `jd guide agents <id>`
- `jd guide commands <name>`

### 출력 내용

- 언제 호출되는지 (트리거 조건)
- 어떻게 활용하는지 (사용 시나리오)
- 적절한 사용법 예시
- 개선점/커스터마이징 제안

### UX

- 로딩 중 Bubbletea 박스 UI로 "Claude Code를 통해 설명 작성 중..." 표시
- AI 응답은 마크다운 형식으로 렌더링
- 캐시된 가이드 표시 시 작성 시간과 재생성 옵션 안내

### 캐싱

- 가이드 저장 경로: `~/.claude/jindo/guides/<type>/<id>.md`
- 캐시된 가이드가 있으면 바로 표시
- `--refresh` (`-r`) 옵션으로 재생성 가능
- 대화형 모드(`-i`)는 캐시 사용 안 함

## Non-goals

- 스킬 수정/맞춤화 (기존 `adapt` 명령 역할)
- 스킬 제작자용 문서 자동 생성
- 여러 스킬 한번에 분석

## Constraints

- 기존 Bubbletea TUI 패턴 재사용 (`internal/tui/`)
- Claude CLI 직접 실행 (`claude -p` 명령 사용)
- 기존 adapt 명령의 프롬프트 관리 패턴 활용 (`internal/prompt/`)

## 구현 파일

| 파일                                       | 역할                   |
| ------------------------------------------ | ---------------------- |
| `internal/cli/guide.go`                    | guide 부모 명령        |
| `internal/cli/guide_skills.go`             | skills 서브커맨드      |
| `internal/cli/guide_hooks.go`              | hooks 서브커맨드       |
| `internal/cli/guide_agents.go`             | agents 서브커맨드      |
| `internal/cli/guide_commands.go`           | commands 서브커맨드    |
| `internal/tui/guide.go`                    | 박스 UI 및 로딩 표시   |
| `internal/guide/guide.go`                  | 가이드 캐싱 로직       |
| `internal/prompt/prompts/guide-skill.md`   | 스킬 안내 프롬프트     |
| `internal/prompt/prompts/guide-hook.md`    | 훅 안내 프롬프트       |
| `internal/prompt/prompts/guide-agent.md`   | 에이전트 안내 프롬프트 |
| `internal/prompt/prompts/guide-command.md` | 명령 안내 프롬프트     |
