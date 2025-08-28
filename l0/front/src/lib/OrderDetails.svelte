<script lang="ts">
    export let orderData: App.Order;

    $: items = orderData?.items ?? [];
    $: total = items.reduce((s: number, i: any) => s + (Number(i.total_price) || 0), 0);
</script>

<aside class="card" aria-label="order card">
    <header class="card-header">
        <div class="customer">
            <div class="label">Customer</div>
            <div class="value">{orderData?.customer_id}</div>
        </div>
        <div class="delivery">
            <div class="label">Delivery</div>
            <div class="value">{orderData?.delivery?.name}</div>
        </div>
    </header>

    <ul class="items">
        {#each items as item (item.rid)}
            <li class="item">
                <div class="left">
                    <div class="name">{item.name}</div>
                    <div class="meta">{item.brand} · {item.size}</div>
                </div>
                <div class="right">{item.total_price} ₽</div>
            </li>
        {/each}
    </ul>

    <div class="total-row">
        <div class="label">Total</div>
        <div class="amount">{total} ₽</div>
    </div>
</aside>

<style>
    .card {
        width: 420px;
        padding: 20px;
        border-radius: 20px;
        box-shadow: 0 10px 30px rgba(16,24,40,0.08);
        box-sizing: border-box;
        font-family: 'Helvetica', sans-serif;
        margin: auto;
    }
    .card-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 20px;
    }
    .label { font-size: 0.8rem; color: #666; }
    .value { font-weight: 600; font-size: 1rem; }

    .items-title { margin: 10px 0; font-size: 1rem; }
    .items { list-style: none; padding: 0; margin: 0 0 20px 0; display: flex; flex-direction: column; gap: 8px; }

    .item { display:flex; justify-content: space-between; align-items: center; padding: 10px 20px; border-radius: 20px; background: #fff; border: 1px solid rgb(200, 200, 200); }
    .left { display:flex; flex-direction: column; }
    .name { font-weight: 600; }
    .meta { font-size: 12px; color: #777; }
    .right { font-weight: 700; }

    .total-row { display:flex; justify-content: space-between; align-items:center; padding: 20px 20px 5px; border-top: 1px solid rgb(200, 200, 200); }
    .amount { font-size: 16px; font-weight: 800; }
</style>
