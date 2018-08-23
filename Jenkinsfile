pipeline {
    agent any
    options { skipDefaultCheckout() }
    stages {
        stage('checkout-normal') {
            when {
                not { branch '**/ready/*' }
            }
            steps  {
                cleanWs()
                checkout scm
                stash name: "repo", includes: "**", useDefaultExcludes: false
            }
        }
        stage('checkout-ready') {
            when {
                branch '**/ready/*'
            }
            steps {
                cleanWs()
                //Using the Pretested integration plugin to checkout out any branch in the ready namespace
                checkout(
                    [$class: 'GitSCM',
                    branches: [[name: '*/ready/**']],
                    doGenerateSubmoduleConfigurations: false,
                    extensions: [[$class: 'CleanBeforeCheckout'],
                    pretestedIntegration(gitIntegrationStrategy: accumulated(),
                    integrationBranch: 'master',
                    repoName: 'origin')],
                    submoduleCfg: [],
                    userRemoteConfigs: [[credentialsId: 'id_rsa', //remember to change credentials and url.
                    url: 'git@github.com:zlepper/welp.git']]])
                stash name: "repo", includes: "**", useDefaultExcludes: false
            }
        }

        stage('build-test') {
            parallel {

                stage('test') {
                    steps {
                        cleanWs()
                        unstash 'repo'
                        sh 'docker run -i --rm -v $PWD:/go/src/github.com/zlepper/welp -w /go/src/github.com/zlepper/welp golang:1.10 /bin/bash -c "go get ./... && go test ./..."'
                    }
                }

                stage('build releases') {
                    when { branch '**/master' }
                    steps {
                        cleanWs()
                        unstash 'repo'
                        sh 'docker run -i --rm -v $PWD:/go/src/github.com/zlepper/welp -w /go/src/github.com/zlepper/welp golang:1.10 /bin/bash -c "go get -d ./... && go run scripts/build/build.go"'
                        stash name: 'artifacts', includes: 'build/**', useDefaultExcludes: false
                    }
                }

                stage('build dockerfile') {
                    agent { label 'docker-releaser' }
                    when {branch '**/master'}
                    steps {
                        unstash 'repo'
                        sh 'docker build -t zlepper/welp:master .'
                    }
                }
            }
        }

        stage('publish') {
            parallel {
                stage('artifacts') {
                    steps {
                        unstash name: 'artifacts'
                        sh 'ls -R'
                        archiveArtifacts 'build/**'
                    }
                }

                stage('docker') {
                    agent { label 'docker-releaser' }
                    when { branch '**/master' }
                    steps {
                        sh 'docker push zlepper/welp:master'
                    }
                }

                stage('pretested') {
                    when {
                        branch '**/ready/*'
                    }
                    steps {
                        //This publishes the commit if the tests have run without errors
                        pretestedIntegrationPublisher()
                    }
                }
            }
        }
    }
}