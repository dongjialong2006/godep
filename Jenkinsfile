pipeline {
    agent any

    stages {
        stage("Prepare") {
            steps {
                echo 'Prepare..'
            }
        }
        stage("Test") {
            steps {
                echo 'Test..'
            }
        }
        stage("Analysis") {
            steps {
                echo 'Analysis..'
            }
        }
        stage('Build') {
            parallel {
                stage('OPS') {
                    agent any
                    steps {
                        echo "On Branch A"
                    }
                }
                stage('SMAC') {
                    agent any
                    steps {
                        sh 'make'
                        archiveArtifacts artifacts: 'bin/*', fingerprint: true
                    }
                }
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}