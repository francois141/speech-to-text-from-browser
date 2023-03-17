#curl -sSL https://get.livekit.io | bash
#curl -sSL https://get.livekit.io/cli | bash
livekit-server --dev
livekit-cli create-token \
    --api-key devkey --api-secret secret \
    --join --room my-first-room --identity user1 \
    --valid-for 24h
#API key: devkey
#API secret: secret

# https://example.livekit.io/