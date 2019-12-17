#SVD LSA example from Alex Thomo's tutorial

#A matrix which will be decomposed.
A = c(1,0,1,0,0, 1,1,0,0,0,0,1,0,0,0,0,1,1,0,0,0,0,0,1,0,0,0,1,1,0,0,0,0,1,0,0,0,0,1,1)
#reshape it into a matrix
dim(A)<-c(5,8) #8 terms and 5 documents
A<-t(A) #needed because we entered the information in row by row


#Now compute the SVD
x<- svd(A) #or svd(A,2,2) to restrict to 2 singular values
sigma <- x$d 
S <- x$u #notation from tutorial is S and U not U and Vt
U <-x$v

#To follow the example, restrict sigma to 2 singular values
#Alternatively, set all the other (smaller singular values to 0, and then get the products so some columns or U and V become 0)
sigma2<-sigma[1:2]
S2<-c(S[,1], S[,2])
dim(S2)<-c(8,2)
#if we had V-Transpose instead, (U here), then would have taken top two rows.
U2<-c(U[,1], U[,2])
dim(U2)<-c(5,2)
U2<-t(U2)

#Now the terms and documents
sigma2Matrix<-c(sigma2[1], 0, 0, sigma2[2])
dim(sigma2Matrix)<-c(2,2)
#sigma2Matrix

#Scaled terms-concept matrix
S2Scaled<-S2 %*% sigma2Matrix
#S2Scaled
rownames(S2Scaled)<-c("romeo", "juliet", "happy", "dagger", "live", "die", "free", "new-hampshire")
colnames(S2Scaled)<-c("LC1", "LC2") #latent concept

#scaled concept-document matrix
U2tScaled<-sigma2Matrix %*% U2
#U2tScaled
colnames(U2tScaled)<-c("Doc1", "Doc2", "Doc3", "Doc4", "Doc5")
rownames(U2tScaled)<-c("LC1", "LC2")


plot(t(U2tScaled), type="p", pch=2, col="blue") #the documents
text(t(U2tScaled), colnames(U2tScaled), pos=3)
lines(S2Scaled, type="p",pch=1, col="grey") #the keywords
text(S2Scaled, rownames(S2Scaled), pos=3)

#Query (dagger+die)
#Find the mid-point of these two and then sort by distance (Cosine used in example)  Cosine is a similarity metric, where vector direction matters, magnitude doesn't.