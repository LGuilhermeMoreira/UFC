package df

import (
	"crypto/rand"
	"errors"
	"math/big"
)

//! Diffie-Hellman

/*
valor tirado de https://datatracker.ietf.org/doc/html/rfc3526,
aparentemente é utilizando como padrão no TLS.

	2048-bit MODP Group

	This group is assigned id 14.

	This prime is: 2^2048 - 2^1984 - 1 + 2^64 * { [2^1918 pi] + 124476 }

	Its hexadecimal value is:

	   FFFFFFFF FFFFFFFF C90FDAA2 2168C234 C4C6628B 80DC1CD1
	   29024E08 8A67CC74 020BBEA6 3B139B22 514A0879 8E3404DD
	   EF9519B3 CD3A431B 302B0A6D F25F1437 4FE1356D 6D51C245
	   E485B576 625E7EC6 F44C42E9 A637ED6B 0BFF5CB6 F406B7ED
	   EE386BFB 5A899FA5 AE9F2411 7C4B1FE6 49286651 ECE45B3D
	   C2007CB8 A163BF05 98DA4836 1C55D39A 69163FA8 FD24CF5F
	   83655D23 DCA3AD96 1C62F356 208552BB 9ED52907 7096966D
	   670C354E 4ABC9804 F1746C08 CA18217C 32905E46 2E36CE3B
	   E39E772C 180E8603 9B2783A2 EC07A28F B5C55DF0 6F4C52C9
	   DE2BCBF6 95581718 3995497C EA956AE5 15D22618 98FA0510
	   15728E5A 8AACAA68 FFFFFFFF FFFFFFFF

	The generator is: 2.
*/
var Prime, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1"+
	"29024E088A67CC74020BBEA63B139B22514A08798E3404DD"+
	"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245"+
	"E485B576625E7EC6F44C42E9A63A3620FFFFFFFFFFFFFFFF", 16)

// numero escolhido foi 2
var Generator = big.NewInt(2)

// gera as chaves secretas. regra: a ∈ [2,P-2]
func GenerateSecretKey(prime *big.Int) (*big.Int, error) {
	min := big.NewInt(2)
	max := new(big.Int).Sub(prime, big.NewInt(4))

	// vai gerar [0,...,max]
	n, err := rand.Int(rand.Reader, max)

	if err != nil {
		return nil, err
	}

	/*
		P = 12
		Numero possiveis: [0,1,2,3,4,5,6,7,8]

		se cair no 8, ainda vai fazer parte da condição

		8 + 2 = 10

		se cair no 0

		0 + 2 = 2

		Ainda faz parte da condição.
	*/
	n.Add(n, min)

	return n, nil
}

// gera as chaves que vão ser compartilhadas
func GenerateSharedKeys(prime, generator, secrectKey *big.Int) (*big.Int, error) {
	n := new(big.Int).Exp(generator, secrectKey, prime)
	if n == nil {
		return nil, errors.New("erro ao gerar chave compartilhada")
	}
	return n, nil
}

func GenerateSharedSecretKey(prime, otherSharedKey, mySecretKey *big.Int) (*big.Int, error) {
	n := new(big.Int).Exp(otherSharedKey, mySecretKey, prime)
	if n == nil {
		return nil, errors.New("erro ao gerar chave secreta compartilhada")
	}
	return n, nil
}
