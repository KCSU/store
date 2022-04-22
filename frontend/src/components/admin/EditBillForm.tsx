import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Bill } from "../../model/Bill";
import { BillDetailsForm } from "./BillDetailsForm";

interface BillProps {
  bill: Bill;
}

export function EditBillForm({ bill }: BillProps) {
  const canWrite = useHasPermission("billing", "write");
  return (
    <BillDetailsForm
      initialValues={bill}
      isDisabled={!canWrite}
      onSubmit={async () => {}}
    >
      Save Changes
    </BillDetailsForm>
  );
}
// const mutation = useEditBill(bill.id);
