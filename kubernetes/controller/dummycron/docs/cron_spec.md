# CronJob Spec

튜토리얼에 따라 몇가지 스펙을 정의한다.

- 스케쥴 (cron in Cronjob)
- 실행한 job에 대한 template (job in CronJob)

몇가지를 더 더하자

- job이 start 하는데에 대한 데드라인
- 한번에 여러가지의 job이 실행될 때 어떻게 해야 할 것인가 에 대한 결정 (wait? stop the old one? run both?)
- Cronjob을 정지할 방법
- old job history에 대한 Limit
    - 우리는 우리 스스로의 status를 절대읽지않기 때문에, job이 실행되었는지 추적할 수 있는 다른 방법이 있어야 한다.
    - old job을 남겨놓음으로서 이를 해결할 수 있을 것

## Controller-tools 사용하기

`controller-tools`는 CRD 설정파일을 생성하는데 사용된다. 이떄, 각 추가적인 metadata 등에 대한 설정은 `// +주석`으로 할 수 있다. 이런 주석은 goDOc에서 각 필드를 설명하는데도 사용될 것이다. 

이와 관련된 코드는 