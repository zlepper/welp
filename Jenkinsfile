pipeline {
    agent none
    options { skipDefaultCheckout() }
    stages {
        stage('checkout-normal') {
            agent {
                node {
                    label 'ubuntu-1'
                }
            }
            when {
                not { branch '**/ready/*' }
            }
            steps  {
                checkout scm
                stash name: "repo", includes: "**", useDefaultExcludes: false
            }
        }
        stage('checkout-ready') {
            agent {
                node {
                    label 'ubuntu-1'
                }
            }
            when {
                branch '**/ready/*'
            }
            steps {
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

        stage('test') {
            agent {
                node {
                    label 'ubuntu-2'
                }
            }
            steps {
                unstash 'repo'
                sh 'docker run -i --rm -v $PWD:/go/src/github.com/zlepper/welp -w /go/src/github.com/zlepper/welp golang:1.10 go get ./... && go test ./...'
            }
        }

        stage('build releases') {
            agent {
                node {
                    label 'ubuntu-3'
                }
            }
            when { branch 'master' }
            steps {
                unstash 'repo'
                sh 'docker run -i --rm -v $PWD:/go/src/github.com/zlepper/welp -w /go/src/github.com/zlepper/welp golang:1.10 go get ./... && go run scripts/build.go'
                stash name: 'artifacts', includes: 'build/**', useDefaultExcludes: false
            }
        }

        stage('publish-artifacts') {
            unstash name: 'artifacts'
            sh 'ls -R'
            archiveArtifacts 'build/**'
        }

        stage('pretested publish') {
            agent {
                node {
                    label 'ubuntu 2'
                }
            }
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