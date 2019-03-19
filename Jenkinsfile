pipeline {
    agent none
    stages {
        stage('Prepare') {
            steps {
                echo 'Prepare..'
            }
        }
        stage('Test') {
            steps {
                echo "Test.."
            }
        }
        stage("Analysis") {
            steps {
                echo 'Analysis..'
            }
        } 
        stage('Dispatch') {
        	parallel {
				stage('OPS') {
                    steps {
                        echo "ops"
                    }
                }
                stage('ZSY') {
                    steps {
                        echo "zsy"
                    }
                }
                stage('SMAC') {
                    stages {
		               stage('Init') {
		                   steps {
		                       echo "Init.."
		                   }
		               }
		               stage('Build') {
		                   steps {
		                       echo "Build.."
		                   }
		               }
		               stage('Deploy') {
		                   steps {
		                       echo "Deploy.."
		                   }
		               }
		            }
                }
			}
        }
    }
}