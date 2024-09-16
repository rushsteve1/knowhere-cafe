// https://passkeys.dev/docs/use-cases/bootstrapping/

async function promptWebauthn() {
  if (
    typeof window.PublicKeyCredential !== "undefined" &&
    typeof window.PublicKeyCredential.isConditionalMediationAvailable ===
      "function"
  ) {
    const available =
      await PublicKeyCredential.isConditionalMediationAvailable();

    if (available) {
      try {
        // Retrieve authentication options for `navigator.credentials.get()`
        // from your server.
        const authOptions = await getAuthenticationOptions();
        // This call to `navigator.credentials.get()` is "set and forget."
        // The Promise will only resolve if the user successfully interacts
        // with the browser's autofill UI to select a passkey.
        const webAuthnResponse = await navigator.credentials.get({
          mediation: "conditional",
          publicKey: {
            ...authOptions,
            // see note about userVerification below
            userVerification: "preferred",
          },
        });
        // Send the response to your server for verification and
        // authenticate the user if the response is valid.
        await verifyAutoFillResponse(webAuthnResponse);
      } catch (err) {
        console.error("Error with conditional UI:", err);
      }
    }
  }
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", () => {
    promptWebauthn(document);
  });
} else {
  processNode(document);
}
