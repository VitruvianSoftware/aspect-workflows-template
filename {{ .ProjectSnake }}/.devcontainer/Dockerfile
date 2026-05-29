# Start with the standard devcontainers base (Ubuntu-based)
# Pinned to 24.04 LTS to avoid the floating ':ubuntu' tag tracking an
# unsupported newer release.
FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

# Install dependencies (direnv + curl for downloading bazelisk)
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
       direnv curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install bazelisk (latest release) — arch-aware (amd64 / arm64)
RUN ARCH="$(dpkg --print-architecture)" \
    && curl -L "https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-${ARCH}" \
      -o /usr/local/bin/bazelisk \
    && chmod +x /usr/local/bin/bazelisk \
    && ln -s /usr/local/bin/bazelisk /usr/local/bin/bazel

# Hook direnv into bash (default Codespaces shell)
RUN echo 'eval "$(direnv hook bash)"' >> /etc/bash.bashrc
