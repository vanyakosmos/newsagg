# newsagg

![build](https://github.com/vanyakosmos/newsagg/actions/workflows/cicd.yml/badge.svg)

Mini news aggregator. It scraps articles from Hacker News and Lobste.rs and push them to 
[News Aggregator](https://telegram.me/newnewsagg) telegram channel.

Aside from that, it also served as a practice ground for using Terraform and GCP.

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

Arch motivation: its free to run.

## Infra and deployment setup

Manual:
- GCP project creation
- github token for terraform
- sentry project and dsn
- telegram channel and bot setup

Infra - run terraform localy with all the secrets:
- artifacts registry
- storage bucket
- service accounts and bindings
- github secrets and vars
- workload identity for github actions

Deployment using terraform via CI/CD:
- cloud run job
- cloud schedule job

CI:
- github action workflow with image build and tf deploy
- connects infra and deploy
