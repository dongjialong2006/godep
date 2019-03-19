pipeline {
    agent any

    stages {
        stage("Prepare") {
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
        stage('Deploy') {
            steps {
            	script {
		        	for(int i = 0; i < 3; i++) {
		        	    stage("Analysis1") {
				            steps {
				                echo 'Analysis1..'
				            }
				        }
						stage("Analysis2") {
				            steps {
				                echo 'Analysis2..'
				            }
				        } 
		            }
		        }
                echo 'Deploy..'
            }
        }
    }
}