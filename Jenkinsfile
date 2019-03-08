pipeline {
    agent any

    stages {
        stage("Prepare") {
            steps {
                echo 'Prepare..'
            }
        }
        stage('Build') {
            steps {
                sh 'make'
                archiveArtifacts artifacts: 'bin/*', fingerprint: true
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

pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh 'make'
                archiveArtifacts artifacts: 'bin/*', fingerprint: true
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