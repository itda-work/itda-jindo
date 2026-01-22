---
number: 36
title: 'feat: jd claudemd guide - CLAUDE.md 베스트 프랙티스 안내'
state: done
labels:
  - feature
assignees: []
created_at: '2026-01-21T17:20:41Z'
updated_at: '2026-01-22T15:30:37Z'
closed_at: '2026-01-22T15:30:37Z'
---

## 개요

CLAUDE.md 작성에 대한 베스트 프랙티스를 AI 기반으로 안내하는 명령 추가

## 기능

- gemini CLI를 활용하여 CLAUDE.md 작성 가이드 제공
- 현재 CLAUDE.md 분석 후 개선점 제안
- 스타일별 템플릿 및 예시 제공

## 사용 예시

```bash
jd claudemd guide                   # 일반 가이드 출력
jd claudemd guide --analyze         # 현재 파일 분석 후 개선점 제안
jd claudemd guide --template        # 템플릿 출력
```

## 제약 사항

- ~~gemini CLI가 설치되어 있어야 함 (optional - 없으면 기본 가이드만 표시)~~
- Claude CLI 사용으로 변경 (기존 guide 명령들과 일관성 유지)

## 구현 내역

### 생성된 파일

- `internal/cli/claudemd_guide.go` - 메인 명령 구현
- `internal/prompt/prompts/guide-claudemd.md` - 프롬프트 템플릿

### 수정된 파일

- `internal/guide/guide.go` - `TypeClaudemd` 추가 (캐시 지원)

### 지원 명령

```bash
jd claudemd guide                   # 일반 가이드 출력
jd claudemd guide --analyze         # 현재 파일 분석 후 개선점 제안
jd claudemd guide --analyze -l      # 로컬 CLAUDE.md 분석
jd claudemd guide --template        # 템플릿 출력
jd claudemd guide -i                # 인터랙티브 모드
jd claudemd guide --format html     # HTML로 브라우저에서 열기
jd claudemd guide --refresh         # 캐시 무시하고 재생성
```

### 기술적 결정

- Claude CLI 사용 (gemini 대신) - 기존 guide skills/commands/hooks/agents와 일관성
- 캐시 지원 (`~/.claude/jindo/guides/claudemd/`) - analyze 모드 제외
- 기존 guide 인프라 재사용 (`guide.RunClaudeWithSpinner`, `guide.PrintGuide` 등)
