pipeline {
    agent any

    stages {
        stage("Prepare") {
            steps {
                echo 'Prepare..'
            }
        }
        stage('Build') {
            parallel {
                stage('BranchA') {
                    agent {
                        label "for-branch-a"
                    }
                    steps {
                        echo "On Branch A"
                    }
                }
                stage('BranchB') {
                    agent {
                        label "for-branch-b"
                    }
                    steps {
                        sh 'make'
                        archiveArtifacts artifacts: 'bin/*', fingerprint: true
                    }
                }
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}