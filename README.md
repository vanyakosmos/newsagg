# newsagg

![build](https://github.com/vanyakosmos/newsagg/actions/workflows/cicd.yml/badge.svg)

My news aggregator

## Project Architecture

```mermaid
graph TD
    sch[Cloud Scheduler]
    job["`Cloud Run Job
    (inside docker container)`"]

    sch -->|runs| job
    job -->|containes| app{Go App}
    app -->|stores tracking data| bucket[(GCS Bucket)]
    app -->|uses bot api| tg((Telegram API))
    app -->|fetch data| hn((HackerNews))
    app -->|fetch data| lo((Lobsters))
```

## Deployemnt

Manual:
- GCP project creation
- github token
- sentry project and dsn
- telegram channel and bot setup

Via terraform:
- artifacts registry
- storage bucket
- service accounts and bindings
- github secrets and vars

Via cd/cd:
- cloud run job
- cloud schedule job
