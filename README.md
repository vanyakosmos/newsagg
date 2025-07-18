# newsagg

![build](https://github.com/vanyakosmos/newsagg/actions/workflows/build.yml/badge.svg)

My news aggregator


## Deployemnt

GCP infra created manually:
- project
- bucket
- service accounts and bindings
- artifacts registry

Github infra created manually:
- secrets

Other:
- sentry

Automatic creation:
- cloud run job
- cloud schedule job


## TODO

- [x] засетапити проєкт на го
- [x] зробити тг бота який зможе періодично робити діставати якусь інфу і слати це як повідомлення в чат чи канал
- [x] зробити скрапер HN і підʼєднати до бота
- [x] зробити скрапер lobste.rs і підʼєднати до бота
- [x] зробити дедуплікатор
- [x] задеплоіти на GCP
- [ ] задеплоіти на GCP використовуючи терраформ (використати terraformer для отримання поточного стану клауда)
- [x] зробити окремий сервісний акк який ранить скеджули і клауд функції
- [x] setup ci
- [x] setup cd

TF migration

- [ ] remove old configs
- [ ] GCP stuff with terraformer(?)
- [ ] auto connect GCP SA to github
