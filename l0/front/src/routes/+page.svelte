<script lang="ts">
    import OrderDetails from '$lib/OrderDetails.svelte'

    let err: string = $state("");
    let inputVal: string = $state("");
    let orderData: App.Order | null = $state(null);

    async function onClick() {
        orderData = null;
        if (inputVal === "") {
            err = "enter order id";
            return;
        }
        err = "";
        const res = await fetch(`http://localhost:8080/orders/${inputVal}`);
        if (!res.ok) {
        const text = await res.text();
            console.error('Backend error:', text);
            err = text;
            return;
        }
        orderData = await res.json();
        console.log(orderData)
    }
</script>

<main>
    <h1>Enter order ID</h1>
    <div class="form">
        <div class="input">
            <input bind:value={inputVal} id="order-input" aria-invalid={err !== ""} class={err ? "error" : ""}/>
            {#if err !== ""}
                <label for="order-input">{err}</label>
            {/if}
        </div>

        <button onclick={onClick}>Get info</button>
    </div>

</main>

{#if orderData}
    <OrderDetails {orderData} />
{/if}

<style>
    main {
        width: 420px;
        margin: auto;
        margin-top: 100px;
        padding: 20px 30px;

        font-family: 'Helvetica', sans-serif;

        display: flex;
        flex-direction: column;
        align-items: center;
    }

    button {
        background-color: rgba(128, 0, 128, 0.8);
        color: white;
        font-size: 16px;
        padding: 14px 24px;
        border-radius: 20px;
        height: 52px;
        width: 100%;

        transition: transform 160ms cubic-bezier(.2,.9,.2,1), box-shadow 160ms ease, background-color 160ms ease;
        transform: translateY(0);
        will-change: transform, box-shadow;
        border: none;
        cursor: pointer;
    }
    button:hover {
        outline: 3px solid rgba(128,0,128,0.22);
        outline-offset: 4px;
    }

    button:focus-visible {  
        outline: 3px solid rgba(128,0,128,0.22);
        outline-offset: 4px;
    }

    h1 {
        font-size: 1.8rem;
        font-weight: 700;
    }

    input {
        border: 2px solid rgb(200, 200, 200);
        border-radius: 20px;
        height: 38px;
        padding: 5px 10px;
        font-size: 1.2rem;
        width: 230px;
    }
    .error {
        border-color: red;  
    }
    label {
        margin-left: 15px;
        color: red;
        font-size: 0.8rem;
    }
    input:focus {
        outline: 3px solid rgba(128,0,128,0.22);
        outline-offset: 3px;
    }
    .form {
        margin: 20px 0;
        width: 100%;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        gap: 20px;
    }
    ul, section {
        width: 100%;
        display: flex;
        flex-direction: column;
        justify-content: start;
    }

    p {
        text-align: left;
        margin: 10px;
    }
</style>
