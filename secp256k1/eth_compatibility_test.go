package secp256k1

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// this function tests whether our library can verify an ECDSA signature
// generated by go-ethereum.
func TestSignature_ShouldVerifyEthereumSignature(t *testing.T) {
	kp := GenerateKeypair()

	msg := [32]byte{1, 2, 3}
	priv, err := kp.private.Encode()
	require.NoError(t, err)
	ethSig, err := secp256k1.Sign(msg[:], priv)
	require.NoError(t, err)

	newSig := &Signature{}
	err = newSig.Decode(ethSig)
	require.NoError(t, err)
	ok, err := kp.Public().Verify(msg[:], newSig)
	require.NoError(t, err)
	assert.True(t, ok)

	pubKey, err := kp.Public().Encode()
	require.NoError(t, err)
	require.Equal(t, 33, len(pubKey))

	ok = secp256k1.VerifySignature(pubKey, msg[:], ethSig[:64])
	require.True(t, ok)
}

// this function tests whether go-ethereum can verify an ECDSA signature
// generated by this library.
func TestSignature_ShouldBeVerifiedByEthereum(t *testing.T) {
	alice := GenerateKeypair()
	oneTime := GenerateKeypair()

	msg := [32]byte{1, 2, 3}
	adaptor, err := alice.AdaptorSign(msg[:], oneTime.public)
	require.NoError(t, err)

	ok, err := alice.Public().VerifyAdaptor(msg[:], oneTime.public, adaptor)
	require.NoError(t, err)
	require.True(t, ok)

	sig, err := adaptor.Decrypt(oneTime.private.key)
	require.NoError(t, err)

	encSig, err := sig.Encode()
	require.NoError(t, err)
	require.Equal(t, 64, len(encSig))

	pubKey, err := alice.Public().Encode()
	require.NoError(t, err)
	require.Equal(t, 33, len(pubKey))

	msg = [32]byte{1, 2, 3}
	verified := secp256k1.VerifySignature(pubKey, msg[:], encSig[:])
	assert.True(t, verified)
}

// test whether go-ethereum can recover a public key from an ECDSA signature
// generated by this library.
func TestSignature_ShouldBeRecoveredByEthereum(t *testing.T) {
	alice := GenerateKeypair()
	oneTime := GenerateKeypair()

	msg := [32]byte{1, 2, 3}
	adaptor, err := alice.AdaptorSign(msg[:], oneTime.public)
	require.NoError(t, err)

	ok, err := alice.Public().VerifyAdaptor(msg[:], oneTime.public, adaptor)
	require.NoError(t, err)
	require.True(t, ok)

	sig, err := adaptor.Decrypt(oneTime.private.key)
	require.NoError(t, err)

	encSig, err := sig.EncodeRecoverable()
	require.NoError(t, err)
	require.Equal(t, 65, len(encSig))

	pubKey, err := alice.Public().Encode()
	require.NoError(t, err)
	require.Equal(t, 33, len(pubKey))

	ethPub, err := secp256k1.RecoverPubkey(msg[:], encSig[:])
	require.NoError(t, err)

	pubKey, err = alice.Public().EncodeDecompressed()
	require.NoError(t, err)

	require.Equal(t, ethPub[1:65], pubKey)
}
