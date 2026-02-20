## コンテナリポジトリの作成

Diarkis クラスタを作成するコンパートメントに対して、コンテナレジストリを作成します

```
# diarkis
COMPARTMENT_ID=ocid1.compartment.oc1..aaaaaaaagywmulolbbttfxurqvwv6bjgo5cytqwidz2niq3paq6zd7lkdvhq
COMPARTMENT_ID={コンパートメントID}
oci artifacts container repository create \
  --compartment-id ${COMPARTMENT_ID} \
  --display-name mars \
  --is--immutable true
```
