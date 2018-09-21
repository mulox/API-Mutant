# API-Mutant (Challenge técnico)

### Verificación de ADN:

Se debe realizar un post al endpoint http://......../mutant  enviando la secuencia de ADN mediante un HTTP POST con un Json el cual tenga el
siguiente formato:

POST → /mutant/
{
“dna”:["ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACTG"]
}

En caso de verificar un mutante, devolvera un HTTP 200-OK, en caso contrario un 403-Forbidden

El metodo ademas de verificar si es mutante o no, verifica que la secuencia de ADN sea Valida y verifica si esta en db. si esta, No continua con la verifica y devuelve el resultado anterior. Si no esta en DB verifica y guarda.

### Estadisticas:

Se debe realizar un get a la siguiente URL: http://...../stats

Será devuelto un JSON con el siguiente formato de ejemplo:

{"count_mutant_dna":4,"count_human_dna":13,"ratio":0.030769}


## Instrucciones de ejecutarlo localmente:

La API esta totalmente desarrollada en GO 1.11, como base de datos se esta utilizando mongoDB.

Como intalar GO, https://golang.org/doc/install

Como instalar MongoDB, 
	[Debian](https://docs.mongodb.com/manual/tutorial/install-mongodb-on-debian/)
	
	Ubuntu
	https://docs.mongodb.com/manual/tutorial/install-mongodb-on-ubuntu/
	Windows
	https://docs.mongodb.com/manual/tutorial/install-mongodb-on-windows/)
	OSx
	https://docs.mongodb.com/manual/tutorial/install-mongodb-on-os-x/)

En el proyecto se incluyeron las siguientes librerías externas, que son necesario instalarlas en el entorno donde se vaya a realizar la prueba:

	* gopkg.in/mgo.v2 (Controlador MongoDB para Go)
	* gopkg.in/mgo.v2/bson (es una implementación de la especificación BSON para Go)
	* github.com/gorilla/mux (enrutador y despachador de URL para Go)

Para poder instalarlos solamente hay que ejecutar las siguientes lineas en tu consola.

	- go get gopkg.in/mgo.v2
	- go get gopkg.in/mgo.v2/bson
	- go get github.com/gorilla/mux

Luego puedes dirigirte al directorio donde se encuentra el proyecto y ejecutar:

	go run *.go


Ademas en el directorio /test se encuentra un archivo php, con el cual se realizo el este y otro archivo dna.php donde había un array con más de 1000 ADNs humanos y mutantes, se dejo solo con 10, para la subida a Github. Este Test se puede ejecutar de la siguiente manera:

	php test.php parallel 10 10

Donde el primer parametro es el tipo de conexion que deseamos realizar (serial|parallel), el segundo es la cantidad de intereaciones y el tercero es la cantidad de conexiones por iteracion que va a realizar.


