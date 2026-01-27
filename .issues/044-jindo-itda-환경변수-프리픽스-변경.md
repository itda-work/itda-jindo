---
number: 44
title: JINDO_ → ITDA_ 환경변수 프리픽스 변경
state: done
labels:
  - enhancement
assignees: []
created_at: '2026-01-27T05:53:25Z'
updated_at: '2026-01-27T06:10:21Z'
closed_at: '2026-01-27T06:10:21Z'
---

## 배경

skills.tts의 API 키 관리를 jindo 통합 설정 시스템으로 이전하면서, 환경변수 프리픽스를 JINDO*에서 ITDA*로 변경한다.
itda-skills 전체에서 사용하는 통합 설정이므로 프리픽스가 ITDA\_가 더 적합하다.

## 변경 내용

- `pkg/config/dotnotation.go`: `toEnvKey()` 함수의 `"JINDO_"` → `"ITDA_"` 변경
- `internal/cli/config_guide.go`: 가이드 내 환경변수 예시 업데이트
- `internal/cli/config_get.go`: env var 표시 관련 코드 업데이트
- 관련 테스트 업데이트

## 완료 조건

- `toEnvKey("common.api_keys.openai")`가 `"ITDA_COMMON_API_KEYS_OPENAI"` 반환
- 기존 `JINDO_` 환경변수 지원 제거 (ITDA\_만 지원)
- 모든 테스트 통과
